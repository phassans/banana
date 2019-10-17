package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/phassans/banana/clients/cloudinary"
	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

func (rtr *router) newListingImageHandler(endpoint postEndpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(hresp{Error: NewAPIError(err)})
			return
		}

		l := shared.Listing{}
		for field, values := range r.Form {
			for _, value := range values {
				if field == "businessId" {
					bid, err := strconv.Atoi(value)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						err = json.NewEncoder(w).Encode(hresp{Error: NewAPIError(err)})
						return
					}
					l.BusinessID = bid
				} else if field == "title" {
					l.Title = value
				} else if field == "discountDescription" {
					l.DiscountDescription = value
				} else if field == "description" {
					l.Description = value
				} else if field == "startDate" {
					l.StartDate = value
				} else if field == "recurringEndDate" {
					l.RecurringEndDate = value
				} else if field == "recurringDays" {
					l.RecurringDays = values
					break
				} else if field == "startTime" {
					l.StartTime = value
				} else if field == "endTime" {
					l.EndTime = value
				}
			}
		}

		m := r.MultipartForm
		images := m.File["images"]

		logger := shared.GetLogger()
		logger = logger.With().
			Str("endpoint", endpoint.GetPath()).Logger()
		logger.Info().Msgf("submit happy hour request")

		if err := ValidateFields(images, r.Form); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(hresp{Error: NewAPIError(err)})
			return
		}

		cloudinaryClient := cloudinary.NewCloudinaryClient(logger)
		fileNames := make([]string, len(images))
		for i, _ := range images {

			uuid, err := uuid.NewRandom()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//for each fileheader, get a handle to the actual file
			file, err := images[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fileName := fmt.Sprintf("upload_images/%s_%s", uuid, images[i].Filename)
			fileNames[i] = fileName
			//create destination file making sure the path is writeable.
			dst, err := os.Create(fileName)
			defer dst.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				err = json.NewEncoder(w).Encode(hresp{Error: NewAPIError(err)})
				return
			}
			//copy the uploaded file to the destination file
			if _, err := io.Copy(dst, file); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				err = json.NewEncoder(w).Encode(hresp{Error: NewAPIError(err)})
				return
			}
		}

		var wg sync.WaitGroup
		wg.Add(len(images))
		for i, _ := range images {
			go func(i int, wg *sync.WaitGroup, l *shared.Listing) {
				defer wg.Done()

				f, err := cloudinaryClient.MustOpen(fileNames[i])
				if err != nil {
					logger.Error().Msgf("error opening file: %s", err)
					return
				}

				values := map[string]io.Reader{
					"file":          f,
					"upload_preset": strings.NewReader(cloudinary.UPLOAD_PRESET),
				}
				cloudinaryResponse, err := cloudinaryClient.Upload(values)
				if err != nil {
					logger.Error().Msgf("error uploading file to cloudinary: %s", err)
					return
				}
				l.ImageLink = cloudinaryResponse.URL
			}(i, &wg, &l)
		}
		wg.Wait()

		listingID, err := rtr.engines.AddListing(&l)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(hresp{Error: NewAPIError(err)})
		}
		fmt.Println(listingID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(hresp{Message: fmt.Sprintf("success added a listing!")})
		return
	}
}

func ValidateFields(images []*multipart.FileHeader, r url.Values) error {
	if len(images) == 0 {
		return helper.ValidationError{Message: fmt.Sprint("submit happy hour failed, missing images!")}
	}

	for field, values := range r {
		for _, value := range values {
			if field == "businessId" {
				bid, err := strconv.Atoi(value)
				if err != nil {
					return err
				}
				if bid == 0 {
					return helper.ValidationError{Message: fmt.Sprint("listing add failed, invalid business ID")}
				}
			} else if field == "title" {
				if value == "" {
					return helper.ValidationError{Message: fmt.Sprint("listing add failed, missing title")}
				}
			} else if field == "discountDescription" {
				if value == "" {
					return helper.ValidationError{Message: fmt.Sprint("listing add failed, missing discountDescription")}
				}
			} else if field == "startDate" {
				if value == "" {
					return helper.ValidationError{Message: fmt.Sprint("listing add  failed, missing startDate")}
				}
			} else if field == "recurringEndDate" {
				if value == "" {
					return helper.ValidationError{Message: fmt.Sprint("listing add  failed, missing recurringEndDate")}
				}
			} else if field == "recurringDays" {
				if len(values) == 0 {
					return helper.ValidationError{Message: fmt.Sprint("listing add  failed, missing recurring days")}
				}
				break
			} else if field == "startTime" {
				if value == "" {
					return helper.ValidationError{Message: fmt.Sprint("listing add  failed, missing startTime")}
				}
			} else if field == "endTime" {
				if value == "" {
					return helper.ValidationError{Message: fmt.Sprint("listing add  failed, missing endTime")}
				}
			}
		}
	}

	return nil
}

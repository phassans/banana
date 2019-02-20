package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/phassans/banana/clients/cloudinary"
	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

const (
	IMAGE_FOLDER_PATH = "upload_images/"
)

type (
	hresp struct {
		Message string    `json:"message,omitempty"`
		Error   *APIError `json:"error,omitempty"`
	}
)

func (rtr *router) newImageHandler(endpoint postEndpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(hresp{Error: NewAPIError(err)})
			return
		}

		m := r.MultipartForm
		images := m.File["images"]
		phoneID := r.FormValue("phoneId")
		name := r.FormValue("name")
		email := r.FormValue("email")
		businessOwner := r.FormValue("businessOwner")
		restaurant := r.FormValue("restaurant")
		city := r.FormValue("city")
		description := r.FormValue("description")

		businessOwnerBool, err := strconv.ParseBool(businessOwner)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(hresp{Error: NewAPIError(err)})
			return
		}

		logger := shared.GetLogger()
		logger = logger.With().
			Str("endpoint", endpoint.GetPath()).
			Str("phoneId", phoneID).
			Str("name", name).
			Str("email", email).
			Bool("businessOwner", businessOwnerBool).
			Str("restaurant", restaurant).
			Str("city", city).
			Str("description", description).Logger()
		logger.Info().Msgf("submit happy hour request")

		if err := Validate(images, restaurant, city); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(hresp{Error: NewAPIError(err)})
			return
		}

		cloudinaryClient := cloudinary.NewCloudinaryClient(logger)
		var cloudinaryResponse cloudinary.Response

		imageLinks := make([]string, len(images))
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

			values := map[string]io.Reader{
				"file":          cloudinaryClient.MustOpen(fileName), // lets assume its this file
				"upload_preset": strings.NewReader(cloudinary.UPLOAD_PRESET),
			}
			cloudinaryResponse, err = cloudinaryClient.Upload(values)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				err = json.NewEncoder(w).Encode(hresp{Error: NewAPIError(err)})
				return
			}

			imageLinks[i] = cloudinaryResponse.URL
		}

		hhID, err := rtr.engines.SubmitHappyHour(phoneID, name, email, businessOwnerBool, restaurant, city, description)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(hresp{Error: NewAPIError(err)})
		}

		for i, _ := range images {
			_, err := rtr.engines.SubmitHappyHourImages(hhID, imageLinks[i])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				err = json.NewEncoder(w).Encode(hresp{Error: NewAPIError(err)})
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(hresp{Message: fmt.Sprintf("success submitted happy hour!")})
		return
	}
}

func Validate(images []*multipart.FileHeader, restaurant string, city string) error {
	if len(images) == 0 {
		return helper.ValidationError{Message: fmt.Sprint("submit happy hour failed, missing images!")}
	}
	if strings.TrimSpace(restaurant) == "" {
		return helper.ValidationError{Message: fmt.Sprint("submit happy hour failed, missing restaurant")}
	}
	if strings.TrimSpace(city) == "" {
		return helper.ValidationError{Message: fmt.Sprint("submit happy hour failed, missing city")}
	}

	return nil
}

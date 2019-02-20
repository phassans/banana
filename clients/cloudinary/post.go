package cloudinary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	BASE_URL      = "https://api.cloudinary.com/v1_1/itshungryhour/image/upload"
	UPLOAD_PRESET = "ouvsftoz"
)

func (c *client) Upload(values map[string]io.Reader) (Response, error) {
	// Prepare a form that you will submit to that URL.
	logger := c.logger
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		var err error
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return Response{}, err
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return Response{}, err
			}
		}
		if _, err := io.Copy(fw, r); err != nil {
			return Response{}, err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", BASE_URL, &b)
	if err != nil {
		return Response{}, err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err", err)
		return Response{}, err
	}
	logger = logger.With().Str("url", BASE_URL).Str("status", resp.Status).Logger()

	if resp.StatusCode != 200 {
		logger = logger.With().Str("body", string(body)).Logger()
		logger.Error().Msgf("doPost non 200 response.json")
		return Response{}, fmt.Errorf("post returned with errorCode: %d", resp.StatusCode)
	}

	// read response.json
	var cloudinaryResponse Response
	err = json.Unmarshal(body, &cloudinaryResponse)
	if err != nil {
		logger = logger.With().Str("error", err.Error()).Logger()
		logger.Error().Msgf("unmarshal error on CrawlResponse")
		return Response{}, err
	}

	fmt.Println(string(body))
	logger.Info().Msgf("doPost success!")
	return cloudinaryResponse, nil
}

func (c *client) MustOpen(f string) (*os.File, error) {
	r, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	return r, nil
}

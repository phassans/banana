package cloudinary

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/phassans/banana/shared"
	"github.com/stretchr/testify/require"
)

func TestClient_Upload(t *testing.T) {
	cloudinaryClient := NewCloudinaryClient(shared.GetLogger())
	f, err := cloudinaryClient.MustOpen("../../upload_images/HungryHourOK_300.jpg")
	require.NoError(t, err)
	//prepare the reader instances to encode
	values := map[string]io.Reader{
		"file":          f, // lets assume its this file
		"upload_preset": strings.NewReader(UPLOAD_PRESET),
	}
	resp, err := cloudinaryClient.Upload(values)
	require.NoError(t, err)
	require.NotNil(t, resp.URL)
	fmt.Println(resp.URL)
}

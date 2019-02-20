package cloudinary

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/phassans/exville/common"
	"github.com/stretchr/testify/require"
)

func TestClient_Upload(t *testing.T) {
	cloudinaryClient := NewCloudinaryClient(common.GetLogger())
	//prepare the reader instances to encode
	values := map[string]io.Reader{
		"file":          cloudinaryClient.MustOpen("../../upload_images/IMG_9614.JPG"), // lets assume its this file
		"upload_preset": strings.NewReader(UPLOAD_PRESET),
	}
	resp, err := cloudinaryClient.Upload(values)
	require.NoError(t, err)
	require.NotNil(t, resp.URL)
	fmt.Println(resp.URL)
}

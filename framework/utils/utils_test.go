package utils_test

import (
	"encoder-project/framework/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsJson(t *testing.T) {
	json := `{
					"id": "12dsd12312-123dasd12-41fsdf1-23121xqc",
					"file_path": "video.mp4",
					"status": "pending"
 			`
	err := utils.IsJson(json)
	require.Nil(t, err)

	json = `acaz`
	err = utils.IsJson(json)
	require.Error(t, err)
}

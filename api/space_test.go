package api

import (
	"testing"

	"github.com/starudream/go-lib/core/v2/utils/testutil"
)

func TestGetSpace(t *testing.T) {
	data, err := GetSpace("x")
	testutil.LogNoErr(t, err, data)
}

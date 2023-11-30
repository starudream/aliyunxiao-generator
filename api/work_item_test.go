package api

import (
	"testing"
	"time"

	"github.com/starudream/go-lib/core/v2/utils/testutil"
)

func TestListWorkItem(t *testing.T) {
	data, err := ListWorkItem(time.Now().Format(time.DateOnly))
	testutil.LogNoErr(t, err, data)
}

func TestListWorkItemTime(t *testing.T) {
	data, err := ListWorkItemTime(time.Now().Format(time.DateOnly), []string{"x", "y"})
	testutil.LogNoErr(t, err, data)
}

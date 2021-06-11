package feed_test

import (
	"testing"

	"github.com/myminicommission/external-data-manager/internal/bsdata/feed"
)

func TestGetAll(t *testing.T) {
	data, err := feed.GetAll()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(data.Entry) == 0 {
		t.Error("data.Entry had no items...")
		t.FailNow()
	}

	for _, entry := range data.Entry {
		t.Log(entry.ID)
	}
}

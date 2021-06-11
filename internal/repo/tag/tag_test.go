package tag_test

import (
	"testing"

	"github.com/myminicommission/external-data-manager/internal/bsdata/feed"
	"github.com/myminicommission/external-data-manager/internal/games/starwars/legion"
	"github.com/myminicommission/external-data-manager/internal/repo/tag"
)

func TestLatest(t *testing.T) {
	minTagValue := "1.7.0"
	repo := legion.RepoName
	data, err := feed.GetAll()
	if err != nil {
		t.Error("could not get feed data", err)
		t.FailNow()
	}

	tagName, err := tag.Latest(repo, &data)
	if err != nil {
		t.Error("could not find tag", err)
		t.FailNow()
	}

	if tagName < minTagValue {
		t.Errorf("tag value was less than expected. Expected: > %s. Actual: %s", minTagValue, tagName)
		t.FailNow()
	}
}

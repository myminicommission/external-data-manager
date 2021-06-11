package tag

import (
	"errors"
	"fmt"
	"strings"

	"github.com/myminicommission/external-data-manager/internal/bsdata/feed"
)

// Latest takes a repo name and a Feed pointer and extracts the highest value tag for the given repo
func Latest(repo string, data *feed.Feed) (string, error) {
	if data == nil {
		return "", errors.New("data is nil")
	}

	tag := "0"
	replaceStr := fmt.Sprintf("https://github.com/BSData/%s/releases/tag/", repo)

	for _, entry := range data.Entry {
		if strings.Contains(entry.ID, fmt.Sprintf("BSData/%s/", repo)) {
			t := strings.Replace(entry.ID, replaceStr, "", 1)
			if t > tag {
				tag = t
			}
		}
	}

	return tag, nil
}

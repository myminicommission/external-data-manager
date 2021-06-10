package wh40k

import (
	"fmt"
	"strings"

	"github.com/myminicommission/external-data-manager/internal/games"
	"github.com/myminicommission/external-data-manager/internal/util"
	"github.com/myminicommission/go-bsdata"
	"github.com/sirupsen/logrus"
)

const (
	gameName = "Warhammer 40,000"
	repoName = "wh40k"
)

// LoadData loads data for Warhammer 40k
func LoadData() ([]games.Mini, error) {
	logrus.Infof("Loading %s data", gameName)

	var minis []games.Mini

	cats, err := bsdata.GetData(repoName)
	if err != nil {
		return nil, err
	}

	for _, cat := range cats {
		for _, e := range cat.EntryLinks.EntryLink {
			if e.Hidden == "false" && !strings.Contains(e.Name, ": ") {
				minis = append(minis, games.Mini{
					Name: fmt.Sprintf("%s - %s", cat.Name, e.Name),
					Game: games.Game{
						Name: gameName,
					},
				})
			}
		}
	}

	return util.RemoveDuplicateMinis(minis), nil
}

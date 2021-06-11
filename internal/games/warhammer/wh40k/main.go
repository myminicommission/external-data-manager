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
	GameName = "Warhammer 40,000"
	RepoName = "wh40k"
)

// LoadData loads data for Warhammer 40k
func LoadData(tag string) ([]games.Mini, error) {
	logrus.Infof("Loading %s data", GameName)

	var minis []games.Mini

	cats, err := bsdata.GetData(RepoName, tag)
	if err != nil {
		return nil, err
	}

	for _, cat := range cats {
		for _, e := range cat.EntryLinks.EntryLink {
			if e.Hidden == "false" && !strings.Contains(e.Name, ": ") {
				minis = append(minis, games.Mini{
					Name: fmt.Sprintf("%s - %s", cat.Name, e.Name),
					Game: games.Game{
						Name: GameName,
					},
				})
			}
		}
	}

	return util.RemoveDuplicateMinis(minis), nil
}

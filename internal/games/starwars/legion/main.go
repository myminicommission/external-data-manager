package legion

import (
	"fmt"
	"strings"

	"github.com/myminicommission/external-data-manager/internal/games"
	"github.com/myminicommission/external-data-manager/internal/util"
	"github.com/myminicommission/go-bsdata"
	"github.com/sirupsen/logrus"
)

const (
	GameName = "Star Wars Legion"
	RepoName = "star-wars-legion"
)

// LoadData loads data for Star Wars Legion
func LoadData(tag string) ([]games.Mini, error) {
	logrus.Infof("Loading %s data", GameName)

	var minis []games.Mini

	cats, err := bsdata.GetData(RepoName, tag)
	if err != nil {
		return minis, err
	}

	for _, cat := range cats {
		for _, e := range cat.EntryLinks.EntryLink {
			e.Name = strings.ReplaceAll(e.Name, "†", "")
			e.Name = strings.TrimSpace(e.Name)

			if e.Hidden == "false" {
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

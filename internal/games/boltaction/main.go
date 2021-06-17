package boltaction

import (
	"fmt"
	"strings"

	"github.com/myminicommission/external-data-manager/internal/games"
	"github.com/myminicommission/external-data-manager/internal/util"
	"github.com/myminicommission/go-bsdata"
	"github.com/sirupsen/logrus"
)

const (
	GameName = "Bolt Action"
	RepoName = "boltaction"
)

// LoadData loads data for Warhammer AoS
func LoadData(tag string) ([]games.Mini, error) {
	logrus.Infof("Loading %s data", GameName)

	var minis []games.Mini

	cats, err := bsdata.GetData(RepoName, tag)
	if err != nil {
		return nil, err
	}

	for _, cat := range cats {
		for _, e := range cat.EntryLinks.EntryLink {
			hidden := e.Hidden != "false"
			hasColon := strings.Contains(e.Name, ": ") // ¯\_(ツ)_/¯
			shouldProcess := !(hidden || hasColon)

			if shouldProcess {
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

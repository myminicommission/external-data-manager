package mmc

import (
	"github.com/myminicommission/external-data-manager/internal/games"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithField("API", "MMC GraphQL")
)

// CreateMini calls the createMini GraphQL mutation for the given mini
func CreateMini(mini games.Mini) (err error) {
	log.WithFields(logrus.Fields{"mini": mini}).Info("creating mini")
	return
}

// GetMini gets the mini from the mmc api
func GetMini(m games.Mini) (mini *games.Mini, err error) {
	log.WithFields(logrus.Fields{"mini": m}).Info("getting mini")
	return
}

// UpdateMini updates the given mini with new data
func UpdateMini(id string, mini games.Mini) (err error) {
	log.WithFields(logrus.Fields{"id": id, "mini": mini}).Info("updating mini")
	return
}

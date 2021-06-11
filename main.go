package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"syscall"

	"github.com/robfig/cron/v3"

	"github.com/myminicommission/external-data-manager/internal/bsdata/feed"
	"github.com/myminicommission/external-data-manager/internal/games"
	"github.com/myminicommission/external-data-manager/internal/games/starwars/legion"
	"github.com/myminicommission/external-data-manager/internal/games/warhammer/wh40k"
	"github.com/myminicommission/external-data-manager/internal/mmc"
	"github.com/sirupsen/logrus"
)

type DataLoader struct {
	Name string
	Fn   func(string) ([]games.Mini, error)
}

var (
	running  = false
	schedule = "@daily"
	feedData *feed.Feed
)

// data loaders
func dataLoaders() []DataLoader {
	return []DataLoader{
		{Name: legion.RepoName, Fn: legion.LoadData},
		{Name: wh40k.RepoName, Fn: wh40k.LoadData},
	}
}

func main() {
	logrus.Info("Starting up")
	if len(os.Args) == 3 {
		schedule = os.Args[2]
	}

	/*
	 * Setup shutdown handler
	 */
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	c := cron.New()
	// process the mini data on the provided schedule
	cID, err := c.AddFunc(schedule, processMiniData)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"cron ID":  cID,
			"schedule": schedule,
		}).Fatal("error starting scheduler")
	}

	logrus.WithFields(logrus.Fields{
		"schedule": schedule,
		"cron ID":  cID,
	}).Info("starting scheduler")

	/*
	* Start the CRON scheduler
	 */
	c.Start()
	<-quit

	logrus.Info("stopping scheduler")
	stopCtx := c.Stop()

	<-stopCtx.Done()

	logrus.Info("application stopped")
}

func processMiniData() {
	if !running {
		running = true

		// load feed data
		if feedData == nil {
			f, err := feed.GetAll()
			if err != nil {
				logrus.Fatal("could not load feed data", err)
			}

			feedData = &f
		}

		for _, loader := range dataLoaders() {
			// capture the function name for logging
			fnLog := newFnLog(&loader)

			fnLog.Info("Executing LoadData")

			tag, err := latestTag(loader.Name)
			if err != nil {
				fnLog.Fatal("could not get latest tag", err)
			}

			minis, err := loader.Fn(tag)
			if err != nil {
				fnLog.Fatal("failed to execute data loder", err)
			}

			for _, mini := range minis {
				handleMini(mini)
			}

			fnLog.Info("LoadData complete")
		}

		// clear feed data
		feedData = nil

		running = false
	} else {
		logrus.Warn("processMiniData already running...")
	}
}

func newFnLog(loader *DataLoader) *logrus.Entry {
	return logrus.WithField("loaderName", loader.Name).WithField(
		"func",
		strings.ReplaceAll(
			runtime.FuncForPC(
				reflect.ValueOf(loader.Fn).Pointer(),
			).Name(),
			"github.com/myminicommission/external-data-manager/internal/games/",
			"",
		),
	)
}

func handleMini(mini games.Mini) {
	log := logrus.WithFields(logrus.Fields{
		"mini": mini,
	})

	log.Info("processing mini")

	// get the existing mini
	existingMini, err := mmc.GetMini(mini)
	if err != nil {
		log.Error("error while calling GetMini", err)
	}
	if err != nil {
		log.Error("error while calling DoesMiniExist", err)
	}

	// was a mini found?
	if existingMini == nil {
		// create the mini
		err := mmc.CreateMini(mini)
		if err != nil {
			log.Error("error while calling CreateMini", err)
		}
	} else {
		// update the mini
		err = mmc.UpdateMini(existingMini.ID, mini)
		if err != nil {
			log.WithFields(logrus.Fields{
				"existingMini": existingMini,
			}).Error("error while calling UpdateMini", err)
		}
	}
}

func latestTag(repo string) (string, error) {
	if feedData == nil {
		return "", errors.New("feedData is nil")
	}

	tag := "0"
	replaceStr := fmt.Sprintf("https://github.com/BSData/%s/releases/tag/", repo)

	for _, entry := range feedData.Entry {
		if strings.Contains(entry.ID, fmt.Sprintf("BSData/%s/", repo)) {
			t := strings.Replace(entry.ID, replaceStr, "", 1)
			if t > tag {
				tag = t
			}
		}
	}

	return tag, nil
}

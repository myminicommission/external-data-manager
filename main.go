package main

import (
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
	"github.com/myminicommission/external-data-manager/internal/games/warhammer/whaos"
	"github.com/myminicommission/external-data-manager/internal/mmc"
	"github.com/myminicommission/external-data-manager/internal/repo/tag"
	"github.com/sirupsen/logrus"

	_ "github.com/joho/godotenv/autoload"
)

type DataLoader struct {
	Name string
	Repo string
	Fn   func(string) ([]games.Mini, error)
}

var (
	running   = false
	schedule  = "@daily"
	feedData  *feed.Feed
	mmcClient mmc.Client
	gameList  []*games.Game
)

// data loaders
func dataLoaders() []DataLoader {
	return []DataLoader{
		{Name: legion.GameName, Repo: legion.RepoName, Fn: legion.LoadData},
		{Name: wh40k.GameName, Repo: wh40k.RepoName, Fn: wh40k.LoadData},
		{Name: whaos.GameName, Repo: whaos.RepoName, Fn: whaos.LoadData},
	}
}

func init() {
	// load the requisite env variables
	url := os.Getenv("MMC_URL")
	if url == "" {
		// if nothing was set then default to the localhost env
		url = "http://localhost:3001/query"
	}
	authToken := os.Getenv("MMC_AUTH_TOKEN") // the defaul value is blank

	// setup the mmc client
	mmcClient = mmc.NewClient(url, authToken)
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

		// load games data
		if len(gameList) == 0 {
			gList, err := mmcClient.ListGames()
			if err != nil {
				logrus.Fatal("could not load games list", err)
			}
			gameList = gList
		}

		for _, loader := range dataLoaders() {
			// capture the function name for logging
			fnLog := newFnLog(&loader)

			fnLog.Info("Executing LoadData")

			// get the game, creating it if neccessary
			game := getGame(loader.Name)
			if game == nil || game.ID == mmc.BlankUUID {
				fnLog.Info("game not found, creating")
				g, err := mmcClient.CreateGame(loader.Name)
				if err != nil {
					fnLog.WithField("error", err).Fatal("could not create game")
				}
				game = g
			}
			fnLog.WithField("game", *game).Info("game loaded")

			tag, err := tag.Latest(loader.Repo, feedData)
			if err != nil {
				fnLog.Fatal("could not get latest tag", err)
			}

			minis, err := loader.Fn(tag)
			if err != nil {
				fnLog.Fatal("failed to execute data loder", err)
			}

			for _, mini := range minis {
				mini.Game = *game
				handleMini(mini)
			}

			fnLog.Info("LoadData complete")
		}

		// clear feed data
		feedData = nil

		// clear game data
		gameList = nil

		running = false
	} else {
		logrus.Warn("processMiniData already running...")
	}
}

func newFnLog(loader *DataLoader) *logrus.Entry {
	return logrus.WithField("loaderName", loader.Name).WithField("loaderRepo", loader.Repo).WithField(
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
	existingMini, err := mmcClient.GetMini(mini)
	if err != nil {
		log.WithField("error", err).Error("error while calling GetMini")
		return
	}

	// was a mini found?
	if existingMini == nil {
		// create the mini
		err := mmcClient.CreateMini(mini)
		if err != nil {
			log.WithField("error", err).Error("error while calling CreateMini")
			return
		}
	} /* else {
		// update the mini
		err = mmcClient.UpdateMini(existingMini.ID, mini)
		if err != nil {
			log.WithFields(logrus.Fields{
				"existingMini": existingMini,
				"error":        err,
			}).Error("error while calling UpdateMini")
			return
		}
	}*/
}

func getGame(name string) *games.Game {
	for _, game := range gameList {
		if game.Name == name {
			return game
		}
	}
	return nil
}

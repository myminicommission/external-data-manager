package mmc

import (
	"context"
	"fmt"

	"github.com/myminicommission/external-data-manager/internal/games"

	"github.com/shurcooL/graphql"
	"github.com/sirupsen/logrus"
)

type Client struct {
	gqlClient *graphql.Client
}

type GameMiniInput struct {
	Game graphql.ID     `json:"game"`
	Name graphql.String `json:"name"`
}

const (
	BlankUUID = "00000000-0000-0000-0000-000000000000"
)

var (
	log = logrus.WithField("API", "MMC GraphQL")
)

// NewClient creates a new Client instance
func NewClient(url, authToken string) Client {
	gqlClient := graphql.NewClient(url, nil)
	return Client{
		gqlClient: gqlClient,
	}
}

// ListGames gets the current list of games
func (c *Client) ListGames() ([]*games.Game, error) {
	var gameList []*games.Game

	var q struct {
		Games []struct {
			ID   graphql.ID
			Name graphql.String
		}
	}

	err := c.gqlClient.Query(context.Background(), &q, nil)
	if err != nil {
		return gameList, err
	}

	for _, game := range q.Games {
		logrus.WithField("game", game).Info("game loaded")
		gameList = append(gameList, &games.Game{
			ID:   fmt.Sprintf("%v", game.ID),
			Name: string(game.Name),
		})
	}

	return gameList, nil
}

// CreateGame creates a new Game through the MMC API
func (c *Client) CreateGame(name string) (*games.Game, error) {
	log := log.WithField("name", name)

	var game *games.Game

	var m struct {
		CreateGame struct {
			ID   graphql.ID
			Name graphql.String
		} `graphql:"createGame(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	log = log.WithField("variables", variables)
	log.Info("Creating new game")

	err := c.gqlClient.Mutate(context.Background(), &m, variables)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to create game")
		return game, err
	}

	game = &games.Game{
		ID:   fmt.Sprintf("%v", m.CreateGame.ID),
		Name: string(m.CreateGame.Name),
	}

	log.WithField("game", game).Info("Game created")

	return game, nil
}

// CreateMini calls the createMini GraphQL mutation for the given mini
func (c *Client) CreateMini(mini games.Mini) error {
	log := log.WithFields(logrus.Fields{"mini": mini})
	log.Info("creating mini")

	var m struct {
		CreateGameMini struct {
			ID   graphql.ID
			Name graphql.String
		} `graphql:"createGameMini(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": GameMiniInput{
			Game: graphql.ID(mini.Game.ID),
			Name: graphql.String(mini.Name),
		},
	}

	err := c.gqlClient.Mutate(context.Background(), &m, variables)
	if err != nil {
		log.WithField("variables", variables).Error("failed to create mini")
		return err
	}

	log.WithField("new mini", m.CreateGameMini).Info("mini created")

	return nil
}

// GetMini gets the mini from the mmc api
func (c *Client) GetMini(m games.Mini) (*games.Mini, error) {
	log.WithFields(logrus.Fields{"mini": m}).Info("getting mini")

	var mini *games.Mini

	var query struct {
		MiniWithName struct {
			ID   graphql.ID
			Name graphql.String
		} `graphql:"miniWithName(name: $name, game: $game)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(m.Name),
		"game": graphql.String(m.Game.Name),
	}

	err := c.gqlClient.Query(context.Background(), &query, variables)
	if err != nil {
		return mini, err
	}

	if query.MiniWithName.ID != BlankUUID {
		mini = &games.Mini{
			ID:   fmt.Sprintf("%v", query.MiniWithName.ID),
			Name: string(query.MiniWithName.Name),
			Game: games.Game{
				Name: m.Game.Name,
			},
		}
	}

	return mini, nil
}

// UpdateMini updates the given mini with new data
func (c *Client) UpdateMini(id string, mini games.Mini) error {
	log.WithFields(logrus.Fields{"id": id, "mini": mini}).Info("updating mini")
	return nil
}

package gameRepo_test

import (
	"testing"

	"github.com/djpiper28/cards-against-humanity/gameLogic"
	"github.com/djpiper28/cards-against-humanity/gameRepo"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	repo := gameRepo.New()
	if repo.GameAgeMap == nil {
		t.Log("The game age map is nil")
		t.FailNow()
	}

	if repo.GameMap == nil {
		t.Log("The game map is nil")
		t.FailNow()
	}

	if repo.GamesByAge == nil {
		t.Log("Games by age is nil")
		t.FailNow()
	}
}

func TestCreateGameFail(t *testing.T) {
	repo := gameRepo.New()
	id, _, err := repo.CreateGame(gameLogic.DefaultGameSettings(), "")
	if err == nil {
		t.Log("When a game errors it should not be made")
		t.FailNow()
	}

	_, found := repo.GameMap[id]
	if found {
		t.Log("The game should not be in the map")
		t.FailNow()
	}

	_, found = repo.GameAgeMap[id]
	if found {
		t.Log("The game should not be in the age map")
		t.FailNow()
	}

	if repo.GamesByAge.Len() > 0 {
		t.Log("The game should not be in the game by age list")
		t.FailNow()
	}
}

func TestCreateGame(t *testing.T) {
	repo := gameRepo.New()

	gameSettings := gameLogic.DefaultGameSettings()
	gameSettings.CardPacks = []*gameLogic.CardPack{{}}
	name := "Dave"
	id, pid, err := repo.CreateGame(gameSettings, name)
	if err != nil {
		t.Log("The game should have been made", err)
		t.FailNow()
	}

	game, _ := repo.GameMap[id]
	if game.PlayersMap[game.GameOwnerId].Name != name {
		t.Log("The player was not made with the correct name")
		t.FailNow()
	}

	assert.Equal(t, pid, game.GameOwnerId, "Game owner should be the returned player ID")

	if game.CreationTime != repo.GameAgeMap[id] {
		t.Log("The age map does not have the game in it")
		t.FailNow()
	}

	if repo.GamesByAge.Front().Value.(gameRepo.GameListPtr) != game {
		t.Log("The games by age list does not contain the game")
		t.FailNow()
	}

	if repo.GamesByAge.Len() != 1 {
		t.Log("The games by age should have length 1")
		t.FailNow()
	}
}

func TestGetGames(t *testing.T) {
	repo := gameRepo.New()

	assert.Equal(t, repo.GetGames(), []*gameLogic.Game{}, "There should be no games in the repo yet")

	gameSettings := gameLogic.DefaultGameSettings()
	gameSettings.CardPacks = []*gameLogic.CardPack{{}}
	gameSettings.CardPacks = []*gameLogic.CardPack{{}}
	name := "Dave"
	id, _, err := repo.CreateGame(gameSettings, name)
	if err != nil {
		t.Log("The game should have been made", err)
		t.FailNow()
	}

	games := repo.GetGames()
	assert.Contains(t, games, repo.GameMap[id], "The game should be in the games returned by the repo")
	assert.Len(t, games, 1, "There should only be one game in the repo")
}

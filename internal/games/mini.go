package games

type Game struct {
	ID   string
	Name string
}

type Mini struct {
	ID   string
	Name string
	Game Game
}

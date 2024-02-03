package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	mapWidth  = 200
	mapHeight = 200
    windowWidth = 51
    windowHeight = 21
    apple = ""
    border = "█"
    snakeHead = ""
    snakeBody = ""
    empty = "⋅"
)

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Every(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type gameState struct {
	length int
	body   [][2]int
	xSpeed int
	ySpeed int
}

func initialState() gameState {
	length := 5
	headX, headY := getRandomStartingPosition(length + 1)
	direction := rand.Intn(4)

	xSpeed := 0
	ySpeed := 0

	switch direction {
	case 0:
		xSpeed = 1
	case 1:
		xSpeed = -1
	case 2:
		ySpeed = 1
	case 3:
		ySpeed = -1
	}

	body := [][2]int{{headX, headY}}
	for i := 1; i < length; i++ {
		body = append(body, [2]int{headX + i*xSpeed, headY + i*ySpeed})
	}

	return gameState{
		length: length,
		xSpeed: xSpeed,
		ySpeed: ySpeed,
		body:   body,
	}
}

func (g *gameState) moveSnake() {
	for i := len(g.body) - 1; i >= 1; i-- {
		g.body[i] = g.body[i-1]
	}
	g.body[0][0] += g.xSpeed
	g.body[0][1] += g.ySpeed
}

func (g *gameState) checkCollisions() bool {
	headX := g.body[0][0]
	headY := g.body[0][1]

	if headX < 0 || headX > mapWidth {
		return true
	}

	if headY < 0 || headY > mapHeight {
		return true
	}

	for i := 1; i < len(g.body); i++ {
		if headX == g.body[i][0] && headY == g.body[i][1] {
			return true
		}
	}

	return false
}

func getRandomStartingPosition(offset int) (int, int) {
	x := rand.Intn(mapWidth-(offset*2)) + offset
	y := rand.Intn(mapWidth-(offset*2)) + offset

	return x, y
}

func (g gameState) Init() tea.Cmd {
	return nil
}

func (g gameState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			g.xSpeed = 0
			g.ySpeed = -1
		case "down", "j":
			g.xSpeed = 0
			g.ySpeed = 1
		case "left", "h":
			g.xSpeed = -1
			g.ySpeed = 0
		case "right", "l":
			g.xSpeed = 1
			g.ySpeed = 0
		case "ctrl+c", "q":
			return g, tea.Quit
		}
	case tickMsg:
		g.moveSnake()

		died := g.checkCollisions()
		if died {
			return g, tea.Quit
		}

		return g, tick()
	}
	return g, nil
}

func (g gameState) View() string {
	return "Hello, world!"
}

func main() {
	p := tea.NewProgram(initialState())
	if _, err := p.Run(); err != nil {
		fmt.Println("Some error happened:", err)
		os.Exit(1)
	}
}

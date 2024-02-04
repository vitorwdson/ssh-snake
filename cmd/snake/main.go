package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	mapWidth     = 200
	mapHeight    = 200
	windowWidth  = 51
	windowHeight = 21
	apple        = "🞜"
	border       = "█"
	snakeHead    = '⦿'
	snakeBody    = '⦾'
	empty        = "⋅"
)

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Every(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type gameState struct {
	dead bool
	length int
	body   [][2]int
	xSpeed int
	ySpeed int
}

func initialState() gameState {
	length := 10
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
		dead: false,
		length: length,
		xSpeed: -xSpeed,
		ySpeed: -ySpeed,
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
	return tick()
}

func (g gameState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if g.ySpeed == 0 {
				g.xSpeed = 0
				g.ySpeed = -1
			}
		case "down", "j":
			if g.ySpeed == 0 {
				g.xSpeed = 0
				g.ySpeed = 1
			}
		case "left", "h":
			if g.xSpeed == 0 {
				g.xSpeed = -1
				g.ySpeed = 0
			}
		case "right", "l":
			if g.xSpeed == 0 {
				g.xSpeed = 1
				g.ySpeed = 0
			}
		case "ctrl+c", "q":
			return g, tea.Quit
		}
	case tickMsg:
		g.moveSnake()

		g.dead = g.checkCollisions()
		if g.dead {
			return g, tea.Quit
		}

		return g, tick()
	}
	return g, nil
}

func buildMap(anchorX, anchorY int) []string {
	snakeMap := make([]string, windowHeight)

	middleX := int(math.Ceil(windowWidth / 2))
	middleY := int(math.Ceil(windowHeight / 2))

	topWall := middleY - anchorY - 1
	rightWall := mapWidth - anchorX + middleX + 1
	bottomWall := mapHeight - anchorY + middleY + 1
	leftWall := middleX - anchorX - 1

	lineStart := max(0, leftWall)
	lineEnd := min(windowWidth, rightWall)
	lineSize := lineEnd - lineStart
	ousideArea := strings.Repeat(" ", lineStart)

	for i := 0; i < windowHeight; i++ {
		if i < topWall || i > bottomWall {
			continue
		}

		if topWall == i || bottomWall == i {
			snakeMap[i] = ousideArea + strings.Repeat(border, lineSize)
			continue
		}

		snakeMap[i] = ousideArea
		if len(snakeMap[i]) > 0 {
			snakeMap[i] += border
		}

		if rightWall > windowWidth {
			snakeMap[i] += strings.Repeat(empty, lineSize)
		} else {
			snakeMap[i] += strings.Repeat(empty, lineSize-1) + border
		}
	}

	return snakeMap
}

func ReplaceCharAt(s string, c rune, i int) string {
    r := []rune(s)
    r[i] = c
    return string(r)
}
func (g gameState) RenderBody(gameMap []string) []string {
	if len(gameMap) != windowHeight {
		return gameMap
	}

	middleX := int(math.Ceil(windowWidth / 2))
	middleY := int(math.Ceil(windowHeight / 2))
	headX := g.body[0][0]
	headY := g.body[0][1]

	offsetX := middleX - headX
	offsetY := middleY - headY

	gameMap[middleY] = ReplaceCharAt(gameMap[middleY], snakeHead, middleX)

	for i := 1; i < len(g.body); i++ {
		x := g.body[i][0] + offsetX
		y := g.body[i][1] + offsetY

		gameMap[y] = ReplaceCharAt(gameMap[y], snakeBody, x)
	}

	return gameMap
}

func (g gameState) View() string {
	if g.dead {
		return ""
	}

	gameMap := buildMap(g.body[0][0], g.body[0][1])
	gameMap = g.RenderBody(gameMap)

	return strings.Join(gameMap, "\n")
}

func main() {
	p := tea.NewProgram(initialState())
	if _, err := p.Run(); err != nil {
		fmt.Println("Some error happened:", err)
		os.Exit(1)
	}

	fmt.Println("Game Over")
}

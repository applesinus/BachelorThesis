package visualizer

import (
	"BachelorThesis/engine/constants"

	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	title = "Bachelor Thesis Visualization"
)

// Game implements ebiten.Game interface.
type Game struct {
	count     int
	algorithm string
	endChan   chan bool
}

// Update proceeds the game state.
func (g *Game) Update() error {
	g.count++

	select {
	case <-g.endChan:
		return ebiten.Termination
	default:
	}

	return nil
}

// Draw draws the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Gray{0x33})
	msg := fmt.Sprintf("TPS: %0.2f\nFrame: %d\nAlgorithm: %s\n", ebiten.ActualTPS(), g.count, g.algorithm)

	ebitenutil.DebugPrint(screen, msg)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constants.WindowWidth, constants.WindowHeight
}

func Start(algorithm string, endChan chan bool) {
	ebiten.SetWindowSize(constants.WindowWidth, constants.WindowHeight)
	ebiten.SetWindowTitle(title)

	game := &Game{}
	game.algorithm = algorithm
	game.endChan = endChan

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

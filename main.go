package main

import (
    "fmt"
    "image/color"
    "math/rand"
    "time"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
    "github.com/hajimehoshi/ebiten/v2/text"
    "golang.org/x/image/font/basicfont"
)

const (
    screenWidth  = 320
    screenHeight = 480
    birdX        = 80
    birdSize     = 20
    pipeWidth    = 40
    pipeGap      = 120
)

type Pipe struct {
    x     float64
    gapY  float64
    scored bool
}

type Game struct {
    birdY   float64
    birdVel float64
    pipes   []Pipe
    score   int
    ticks   int
}

func NewGame() *Game {
    rand.Seed(time.Now().UnixNano())
    return &Game{birdY: screenHeight / 2}
}

func (g *Game) Reset() {
    g.birdY = screenHeight / 2
    g.birdVel = 0
    g.pipes = nil
    g.score = 0
    g.ticks = 0
}

func (g *Game) Update() error {
    if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
        g.birdVel = -6
    }
    g.birdVel += 0.4
    g.birdY += g.birdVel

    if g.ticks%90 == 0 {
        gapY := 60 + rand.Float64()*(screenHeight-2*60-pipeGap)
        g.pipes = append(g.pipes, Pipe{x: screenWidth, gapY: gapY})
    }
    g.ticks++

    for i := range g.pipes {
        g.pipes[i].x -= 2
    }

    for len(g.pipes) > 0 && g.pipes[0].x+pipeWidth < birdX {
        if !g.pipes[0].scored {
            g.score++
        }
        g.pipes = g.pipes[1:]
    }

    if g.birdY+birdSize > screenHeight || g.birdY < 0 {
        g.Reset()
    }

    for _, p := range g.pipes {
        if birdX+birdSize > p.x && birdX < p.x+pipeWidth {
            if g.birdY < p.gapY || g.birdY+birdSize > p.gapY+pipeGap {
                g.Reset()
            }
        }
    }

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    screen.Fill(color.RGBA{135, 206, 235, 255})

    ebitenutil.DrawRect(screen, birdX, g.birdY, birdSize, birdSize, color.RGBA{255, 255, 0, 255})

    for _, p := range g.pipes {
        ebitenutil.DrawRect(screen, p.x, 0, pipeWidth, p.gapY, color.RGBA{0, 128, 0, 255})
        ebitenutil.DrawRect(screen, p.x, p.gapY+pipeGap, pipeWidth, screenHeight-(p.gapY+pipeGap), color.RGBA{0, 128, 0, 255})
    }

    text.Draw(screen, fmt.Sprintf("Score: %d", g.score), basicfont.Face7x13, 10, 20, color.Black)
}

func (g *Game) Layout(outW, outH int) (int, int) {
    return screenWidth, screenHeight
}

func main() {
    ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
    ebiten.SetWindowTitle("Flappy Bird")
    if err := ebiten.RunGame(NewGame()); err != nil {
        panic(err)
    }
}


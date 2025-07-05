package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Material int

const (
	worldWidth             = 200
	worldHeight            = 200
	MaterialEmpty Material = iota
	MaterialSand
	MaterialDirt
)

type Cell struct {
	Material Material
}

var world [worldWidth][worldHeight]Cell
var frameCounter int

type Game struct{}

func (g *Game) Update() error {
	world = GenerateGround(world, 100)

	if frameCounter > 50 && frameCounter < 100 {
		// place falling sand at the top
		for x := 50; x < 100; x++ {
			//log.Println("framecounter: ", frameCounter)
			world[x][0].Material = MaterialSand
		}
	}

	if (frameCounter > 200 && frameCounter < 250) && frameCounter%2 == 0 {
		// place falling sand at the top
		for x := range worldWidth {
			if x%2 == 0 {
				//log.Println("framecounter: ", frameCounter)
				world[x][0].Material = MaterialSand
			}
		}
	}

	// simple sand simulation
	for y := worldHeight - 2; y >= 0; y-- {
		for x := 0; x < worldWidth; x++ {
			if world[x][y].Material == MaterialSand && world[x][y+1].Material == MaterialEmpty {
				//				log.Printf("Pixel at world[%d][%d] has Material = %d\n", x, y, world[x][y].Material)
				world[x][y+1].Material = MaterialSand
				world[x][y].Material = MaterialEmpty
			}
		}
	}
	frameCounter++
	return nil
}

func GenerateGround(world [worldWidth][worldHeight]Cell, rows int) [worldWidth][worldHeight]Cell {
	low := 0
	high := len(world)

	if frameCounter == 0 {

		for i := range world {
			if i >= len(world)-rows {
				log.Printf("we made it into the zone")

				for j := range world {
                    log.Printf("i: %d\nj: %d\nlow: %d\nhigh: %d\n\n", i, j, low, high)
                    //							log.Printf("%d + %d = %d", i, j, i+j)
					if j >= low && j <= high {
                        world[j][i].Material = MaterialDirt
					}
                    low = low+1
                    high = high-1
				}
			}
			//world[i][worldHeight-1].Material = MaterialDirt
			//world[i][worldHeight-2].Material = MaterialDirt
			//world[i][worldHeight-3].Material = MaterialDirt
			//world[i][worldHeight-4].Material = MaterialDirt
			//world[i][worldHeight-5].Material = MaterialDirt
		}
	}

	return world
}

func (g *Game) Draw(screen *ebiten.Image) {
	// clear
	screen.Fill(color.RGBA{255, 0, 0, 255})

	// draw world
	for y := 0; y < worldHeight; y++ {
		for x := 0; x < worldWidth; x++ {
			switch world[x][y].Material {
			case MaterialSand:
				world[x][y].Material = MaterialSand
				screen.Set(x, y, color.RGBA{194, 178, 128, 255})
			case MaterialDirt:
				world[x][y].Material = MaterialDirt
				screen.Set(x, y, color.RGBA{123, 63, 0, 255})
			case MaterialEmpty:
			default:
				world[x][y].Material = MaterialEmpty
			}
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return worldWidth, worldHeight
}

func main() {
	ebiten.SetWindowSize(worldWidth*3, worldHeight*3)
	ebiten.SetWindowTitle("Sand Simulation")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

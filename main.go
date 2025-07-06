package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Material int

const (
	worldWidth             = 50
	worldHeight            = 50
	MaterialEmpty Material = iota
	MaterialSand
	MaterialDirt
)

type Cell struct {
	Material Material
}

type Player struct {
	Position []int
	HP       int
}

var world [worldWidth][worldHeight]Cell
var player Player
var frameCounter int
var jumpCooldown = true
var jumpWindow = true

type Game struct{}

func (g *Game) Update() error {
	//log.Printf("Player x position: %d\n", player.Position[0])
	//log.Printf("Player y position: %d\n", player.Position[1])

	//player gravity should happen even if the player is pressing space
	if world[player.Position[0]][player.Position[1]+1].Material == MaterialEmpty {
		log.Printf("In the air: %d\n", world[player.Position[0]][player.Position[1]+1].Material)
		if jumpWindow == false {
			player.Position[1] += 1
		} else {
            player.Position[1] += 1
        }
	} else {
		if player.Position[1]+1 < worldHeight {
			log.Printf("Material under: %d\n", world[player.Position[0]][player.Position[1]+1].Material)
		} else {
			log.Printf("No material under\n")
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		if player.Position[0] <= 0 {
		} else {
			if world[player.Position[0]-1][player.Position[1]].Material == MaterialEmpty {
				player.Position[0] -= 1
			} else if world[player.Position[0]-1][player.Position[1]-1].Material == MaterialEmpty {
				player.Position[1] -= 1
				player.Position[0] -= 1
			}
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		if player.Position[0]+1 >= worldWidth {
		} else {
			if world[player.Position[0]+1][player.Position[1]].Material == MaterialEmpty {
				player.Position[0] += 1
			} else if world[player.Position[0]+1][player.Position[1]-1].Material == MaterialEmpty {
				player.Position[1] -= 1
				player.Position[0] += 1
			}
		}
	}

	jumpResetTimer := time.NewTimer(2 * time.Second)
	jumpWindowTimer := time.NewTimer(1 * time.Second)

	//log.Printf("Player position: x:%d, y:%d | \n", player.Position[0], player.Position[1])
	//log.Printf("Current material: %d\n", world[player.Position[0]][player.Position[1]].Material)
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		log.Printf("Spacebar\n")
		log.Println("jumpwindow: ", jumpWindow)

        if world[player.Position[0]][player.Position[1]+1].Material != MaterialEmpty {
            go func() {
                log.Printf("jump reset timer started.")
                <-jumpResetTimer.C
                log.Printf("%s\n", "jump reset timer has fired.")
                jumpCooldown = true
            }()

            //log.Printf("Material Empty: %d\n", MaterialEmpty)
            //log.Printf("world[player.Position[0]][player.Position[1]].Material: %d\n", world[player.Position[0]][player.Position[1]].Material)
            //log.Printf("world[player.Position[0]][player.Position[1]+1].Material: %d\n", world[player.Position[0]][player.Position[1]+1].Material)
            go func() {
                log.Printf("\n\njump window timer started.\n\n")
                <-jumpWindowTimer.C
                log.Printf("jump window timer has fired.")
                jumpWindow = false 
            }()
        }

		if jumpWindow == true && jumpCooldown == true {
			log.Printf("\n\nJUMPING????\n\n")
			if player.Position[1] > 0 && player.Position[1] < worldWidth {
				player.Position[1] -= 1
			}
		} else if jumpCooldown == true {
			log.Printf("\n\nJUMPING OFF GROUND\n\n")
			if world[player.Position[0]][player.Position[1]+1].Material != MaterialEmpty {
				log.Printf("Jump!\n")
				if player.Position[1] > 0 && player.Position[1] < worldWidth {
					player.Position[1] -= 1
				}
			}
		} else {

		}
		jumpCooldown = false
	}
	//if ebiten.IsKeyPressed(ebiten.KeyUp) {
	//	if player.Position[1] <= 0 {
	//	} else {
	//		player.Position[1] -= 1
	//	}
	//} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
	//	if player.Position[1] >= worldHeight-1 {
	//	} else {
	//		player.Position[1] += 1
	//	}
	//}

	if frameCounter > 0 && frameCounter < 10 {
		// place falling sand at the top
		for x := range worldWidth {
			if x > 20 && x < 30 {
				//log.Printf("Generating sand at world[%d][%d]", x, 0)
				//log.Println("framecounter: ", frameCounter)
				world[x][0].Material = MaterialSand
			}
		}
	}

	//if (frameCounter > 200 && frameCounter < 250) && frameCounter%2 == 0 {
	//	// place falling sand at the top
	//	for x := range worldWidth {
	//		for range 5 {
	//			//log.Println("framecounter: ", frameCounter)
	//			world[x][0].Material = MaterialSand
	//		}
	//	}
	//}

	// simple sand simulation
	for x := range worldHeight {
		for y := worldWidth - 1; y >= 0; y-- {
			if y+1 <= worldHeight-1 {
				if world[x][y].Material == MaterialSand {
					if world[x][y+1].Material == MaterialEmpty {
						//log.Printf("Pixel at world[%d][%d] has Material = %d\n", x, y, world[x][y].Material)
						world[x][y].Material = MaterialEmpty
						world[x][y+1].Material = MaterialSand
					} else if world[x-1][y+1].Material == MaterialEmpty {
						world[x][y].Material = MaterialEmpty
						world[x-1][y+1].Material = MaterialSand
					} else if world[x+1][y+1].Material == MaterialEmpty {
						world[x][y].Material = MaterialEmpty
						world[x+1][y+1].Material = MaterialSand
					}
				}
			}
		}
	}

	frameCounter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// clear
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// draw world
	for x := range worldHeight {
		for y := range worldWidth {
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
				screen.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}
	screen.Set(player.Position[0], player.Position[1], color.RGBA{255, 0, 0, 255})
	screen.Set(player.Position[0], player.Position[1]+1, color.RGBA{0, 255, 0, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return worldWidth, worldHeight
}

func main() {
	var exampleworld [worldHeight][worldWidth]Cell
	for i := range worldHeight {
		for j := range worldWidth {
			if j > 41 {
				exampleworld[i][j] = Cell{Material: MaterialDirt}
			} else {
				exampleworld[i][j] = Cell{Material: MaterialEmpty}
			}
		}
	}

	bytes, err := json.Marshal(exampleworld)
	if err != nil {
		log.Fatal("error marshalling example world:", err)
	}
	err = os.WriteFile("exampleworld.json", bytes, 666)
	if err != nil {
		log.Fatal(err)
	}

	world, err = NewWorld()
	if err != nil {
		log.Fatal("error generating world:", err)
	}
	fmt.Println(world)

	player = Player{
		Position: []int{45, 1},
		HP:       100,
	}

	ebiten.SetWindowSize(worldWidth*10, worldHeight*10)
	ebiten.SetWindowTitle("Sand Simulation")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func NewWorld() ([worldHeight][worldWidth]Cell, error) {
	var newWorld [worldHeight][worldWidth]Cell
	bytes, err := os.ReadFile("exampleworld.json")
	if err != nil {
		return newWorld, err
	}

	err = json.Unmarshal(bytes, &newWorld)
	if err != nil {
		return newWorld, err
	}

	return newWorld, nil
}

package main

import (
	"consoleEngine/Engine2D"
	GUI "consoleEngine/gui"
	"image"
	"image/png"
	"log"
	"os"
	"time"
)

var laser image.Image

func main() {

	// shrek, err := loadTexture("textures/shrek.png")
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	laser, err := loadTexture("../../textures/spongebob.png")
	if err != nil {
		log.Println(err)
		return
	}

	gui := GUI.New()
	screen := Engine2D.NewScreen(gui)

	shrekEntity := Engine2D.NewPhysicsEntity(
		screen,
		laser,
		GUI.Coords{X: 0, Y: 0},
		GUI.Coords{X: screen.ScreenSize.X, Y: screen.ScreenSize.Y},
		0,
		nil, //handleInput,
		false,
		GUI.Coords{X: 0, Y: 0},
		100,
	)

	// shrekEntity2 := Engine2D.NewPhysicsEntity(
	// 	screen,
	// 	shrek,
	// 	GUI.Coords{X: 100, Y: 0},
	// 	GUI.Coords{X: 25, Y: 15},
	// 	1,
	// 	nil,
	// 	false,
	// 	GUI.Coords{X: -1, Y: 0},
	// 	50)

	// shrekEntity3 := Engine2D.NewPhysicsEntity(
	// 	screen,
	// 	shrek,
	// 	GUI.Coords{X: 200, Y: 0},
	// 	GUI.Coords{X: 25, Y: 15},
	// 	2,
	// 	nil,
	// 	false,
	// 	GUI.Coords{X: 1, Y: 0},
	// 	1)

	screen.AddEntity(shrekEntity)
	// screen.AddEntity(shrekEntity2)
	// screen.AddEntity(shrekEntity3)

	//	go Screensaver(screen, shrekEntity)
	//	go Screensaver(screen, shrekEntity2)
	//	go Screensaver(screen, shrekEntity3)

	ticker := time.NewTicker(1000 / 140 * time.Millisecond)
	for {
		screen.Paint()
		<-ticker.C
	}
}

func loadTexture(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return png.Decode(f)
}

/*
func Screensaver(screen *GUI.Screen, entity GUI.IEntity) {
	nextStepX := 1.0
	nextStepY := 1.0
	rotationStep := -3.0
	ticker := time.NewTicker(1000 / 30 * time.Millisecond)
	for {
		entity.Pos.X += nextStepX
		entity.Pos.Y += nextStepY
		entity.Rotation += rotationStep

		if entity.Pos.X+entity.Size.X >= screen.ScreenSize.X || entity.Pos.X < 0 {
			nextStepX = -nextStepX
			rotationStep = -rotationStep
		}

		if entity.Pos.Y+entity.Size.Y >= screen.ScreenSize.Y || entity.Pos.Y < 0 {
			nextStepY = -nextStepY
		}
		<-ticker.C
	}
}
*/

func handleInput(entity Engine2D.IPhysicsEntity, key Engine2D.KeyCode) {
	step := 3.0
	pos := entity.GetPosition()
	mom := entity.GetMomentum()
	rot := entity.GetRotation()
	screen := entity.GetScreen()

	switch key {
	case Engine2D.KeyCode_W:
		if pos.Y+entity.GetSize().Y >= screen.ScreenSize.Y {
			mom.Y += -5
		}
	case Engine2D.KeyCode_A:
		mom.X -= step
	case Engine2D.KeyCode_S:
	case Engine2D.KeyCode_D:
		mom.X += step
	case Engine2D.KeyCode_Space:
		laserEntity := Engine2D.NewEntity(
			entity.GetScreen(),
			laser,
			GUI.Coords{X: 0, Y: 0},
			GUI.Coords{X: 10, Y: 10},
			3,
			nil,
		)

		laserEntity.Pos = pos
		laserEntity.Visible = true
		screen.AddEntity(laserEntity)
		stepX := 1
		stepY := 0
		go laserMovement(screen, laserEntity, stepX, stepY)
	}
	entity.SetMomentum(mom)
	entity.SetRotation(rot)

}

func laserMovement(screen *Engine2D.Screen, entity Engine2D.Entity, stepX int, stepY int) {
	pos := entity.GetPosition()
	for pos.X < screen.ScreenSize.X || pos.Y < screen.ScreenSize.Y || pos.X > 0 || pos.Y > 0 {
		pos.X += float64(stepX)
		pos.Y += float64(stepY)
		entity.SetPosition(pos)
		time.Sleep(10 * time.Millisecond)
	}

	screen.RemoveEntity(entity)
}

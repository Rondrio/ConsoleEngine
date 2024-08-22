package main

import (
	GUI "consoleEngine/gui"
	"image"
	"image/png"
	"log"
	"os"
	"time"
)

var laser image.Image

func main() {

	shrek, err := loadTexture("textures/shrek.png")
	if err != nil {
		log.Println(err)
		return
	}

	laser, err = loadTexture("textures/laser.png")
	if err != nil {
		log.Println(err)
		return
	}

	gui := GUI.New()
	screen := GUI.NewScreen(gui)

	shrekEntity := GUI.NewPhysicsEntity(
		screen,
		shrek,
		GUI.Coords{X: 0, Y: 0},
		GUI.Coords{X: 25, Y: 15},
		0,
		handleInput,
		false,
		GUI.Coords{X: 0, Y: 0},
		100,
	)

	shrekEntity2 := GUI.NewPhysicsEntity(
		screen,
		shrek,
		GUI.Coords{X: 100, Y: 0},
		GUI.Coords{X: 25, Y: 15},
		1,
		nil,
		false,
		GUI.Coords{X: -1, Y: 0},
		50)

	shrekEntity3 := GUI.NewPhysicsEntity(
		screen,
		shrek,
		GUI.Coords{X: 200, Y: 0},
		GUI.Coords{X: 25, Y: 15},
		2,
		nil,
		false,
		GUI.Coords{X: 1, Y: 0},
		1)

	screen.AddEntity(shrekEntity)
	screen.AddEntity(shrekEntity2)
	screen.AddEntity(shrekEntity3)

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

func handleInput(entity GUI.IPhysicsEntity, key GUI.KeyCode) {
	step := 3.0
	pos := entity.GetPosition()
	mom := entity.GetMomentum()
	rot := entity.GetRotation()
	screen := entity.GetScreen()

	switch key {
	case GUI.KeyCode_W:
		if pos.Y+entity.GetSize().Y >= screen.ScreenSize.Y {
			mom.Y += -5
		}
	case GUI.KeyCode_A:
		mom.X -= step
	case GUI.KeyCode_S:
	case GUI.KeyCode_D:
		mom.X += step
	case GUI.KeyCode_Space:
		laserEntity := GUI.NewEntity(
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

func laserMovement(screen *GUI.Screen, entity GUI.Entity, stepX int, stepY int) {
	pos := entity.GetPosition()
	for pos.X < screen.ScreenSize.X || pos.Y < screen.ScreenSize.Y || pos.X > 0 || pos.Y > 0 {
		pos.X += float64(stepX)
		pos.Y += float64(stepY)
		entity.SetPosition(pos)
		time.Sleep(10 * time.Millisecond)
	}

	screen.RemoveEntity(entity)
}

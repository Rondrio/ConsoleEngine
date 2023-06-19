package main

import (
	GUI "consoleEngine/gui"
	"image"
	"image/png"
	"log"
	"os"
	"time"
)

func main() {

	spongeBoss, err := loadTexture("textures/shrek.png")
	if err != nil {
		log.Println(err)
		return
	}

	grid, err := loadTexture("textures/grid.png")
	if err != nil {
		log.Println(err)
		return
	}
	GUI := GUI.New()
	GUI.PaintImage(spongeBoss, 0, 0, 70, 150)
	GUI.PaintImage(grid, 0, 0, 70, 150)
	for {
		GUI.PrintBuffer()
		time.Sleep(1000 / 60 * time.Millisecond)
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

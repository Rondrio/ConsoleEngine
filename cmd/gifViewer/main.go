package main

import (
	"bufio"
	GUI "consoleEngine/gui"
	"image/gif"

	"os"
	"time"
)

func main() {
	gif, err := loadImage("../../textures/maxwell-spin.gif")
	if err != nil {
		panic(err)
	}
	gui := GUI.New()

	for {
		for i, img := range gif.Image {
			// continue
			gui.PaintImage(img, GUI.Coords{0, 0}, gui.GetScreenSize())
			gui.PrintBuffer()
			time.Sleep(time.Duration(gif.Delay[i]) * (time.Second / 100))
		}
	}

}

func loadImage(path string) (*gif.GIF, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := bufio.NewReader(file)

	return gif.DecodeAll(buffer)
}

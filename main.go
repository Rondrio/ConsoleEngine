package main

import (
	"bytes"
	"consoleEngine/GUI"
	"image/gif"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	file, err := os.Open("testgif.gif")
	if err != nil {
		log.Println("failed to open file", err)
		return
	}
	defer file.Close()

	g, err := gif.DecodeAll(file)
	if err != nil {
		log.Println("failed to decode gif", err)
		return
	}

	//ticker := time.NewTicker(33 * time.Millisecond)
	var b bytes.Buffer

	frameChan := make(chan int)
	paintChan := make(chan int)

	go func(gif *gif.GIF,b io.Writer,frameChan,paintChan chan int) {
			for _,image := range g.Image {
				select {
				case i := <-paintChan:
					GUI.PaintGif(image, b)
					frameChan <- i
				}
			}
	}(g,&b,frameChan,paintChan)

	paintChan <- 1
	ticker := time.NewTicker(33*time.Millisecond)
	for  {
		var copy bytes.Buffer
		select {
			case i:=<-frameChan:
				copy = b
				i+=1
				paintChan <- i

			case <-ticker.C:
				//GUI.PaintGif(image, &b)
				io.Copy(os.Stdout, &copy)
			//	time.Sleep(1*time.Second)
			}
		//	<-ticker.C
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

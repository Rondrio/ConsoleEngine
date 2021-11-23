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

	/*	//ticker := time.NewTicker(33 * time.Millisecond)
		var b bytes.Buffer

		frameChan := make(chan int)
		paintChan := make(chan int)

		func(gif *gif.GIF, b io.Writer, frameChan, paintChan chan int) {
			for _, image := range g.Image {
				//	select {
				//	case i := <-paintChan:
				GUI.PaintGif(image, b)
				//		frameChan <- i
				io.Copy(os.Stdout, b)
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()

				//	}
			}
		}(g, &b, frameChan, paintChan)
		/*
		   	paintChan <- 0
		   	ticker := time.NewTicker(33 * time.Millisecond)
		   	for {
		   		var copy bytes.Buffer
		   		select {
		   		case i := <-frameChan:
		   			copy = b
		   			i += 1

		   		case <-ticker.C:
		   			//GUI.PaintGif(image, &b)
		   			io.Copy(os.Stdout, &copy)
		   			paintChan <- 0
		   			//	time.Sleep(1*time.Second)
		   		}
		   		cmd := exec.Command("clear")
		   		cmd.Stdout = os.Stdout
		   		cmd.Run()
		   	}
		   }

	*/

	bufferChan := make(chan bytes.Buffer)
	framePaintedChan := make(chan int)

	go WriteFrame(g, bufferChan, framePaintedChan)
	go ScreenClearer(bufferChan)

	//	var b bytes.Buffer
	for { //_, image := range g.Image {
		b := <-bufferChan
		//GUI.PaintGif(image, &b)
		io.Copy(os.Stdout, &b)
	}
}

func WriteFrame(gif *gif.GIF, bufferChan chan bytes.Buffer, framePaintedChan chan int) {
	ticker := time.NewTicker(33 * time.Millisecond)
	for _, image := range gif.Image {
		<-ticker.C
		var b bytes.Buffer
		GUI.PaintGif(image, &b)
		bufferChan <- b
	}
}

func ScreenClearer(bufferChan chan bytes.Buffer) {
	for {
		<-bufferChan
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

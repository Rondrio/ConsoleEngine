package gui

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"os"

	paintColor "github.com/fatih/color"
	"golang.org/x/term"
)

const (
	location_foreground = location("foreground")
	location_background = location("background")
)

var (
	COLORS_BACKGROUND_4bit = map[byte]paintColor.Attribute{
		/* A12830 - FF584F*/
		0b000: paintColor.BgHiBlack,
		//{R: 0x29, G: 0x29, B: 0x29, A: 0xFF}: paintColor.BgHiBlack, // Gray M

		0b001: paintColor.BgHiBlue,
		//{R: 0x99, G: 0xCC, B: 0xFF, A: 0xFF}: paintColor.BgHiBlue,

		0b010: paintColor.BgHiGreen,
		//0b010101: paintColor.BgHiGreen,

		0b011: paintColor.BgHiCyan,
		//{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}: paintColor.BgHiCyan,

		0b100: paintColor.BgHiRed, // 	Red Range
		//{R: 0x99, G: 0x00, B: 0x00, A: 0xFF}: paintColor.BgHiRed, // Red Range

		0b101: paintColor.BgHiMagenta, // Magenta M
		//{R: 0x99, G: 0x00, B: 0x99, A: 0xFF}: paintColor.BgHiMagenta, // Magenta M

		0b110: paintColor.BgHiYellow, // Orange
		//		{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}: paintColor.BgHiYellow, //paintColor.FgHiGreen,

		0b111: paintColor.BgHiWhite, // light gray
		//{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}: paintColor.BgHiWhite,
	}

	COLORS_FOREGROUND_4bit = map[byte]paintColor.Attribute{
		/* A12830 - FF584F*/
		/* A12830 - FF584F*/
		0b000: paintColor.FgHiBlack,
		//{R: 0x29, G: 0x29, B: 0x29, A: 0xFF}: paintColor.BgHiBlack, // Gray M

		0b001: paintColor.FgHiBlue,
		//{R: 0x99, G: 0xCC, B: 0xFF, A: 0xFF}: paintColor.BgHiBlue,

		0b010: paintColor.FgHiGreen,
		//	{R: 0, G: 1, B: 0, A: 0xFF}: paintColor.BgHiGreen,

		0b011: paintColor.FgHiCyan,
		//{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}: paintColor.BgHiCyan,

		0b100: paintColor.FgHiRed, // 	Red Range
		//{R: 0x99, G: 0x00, B: 0x00, A: 0xFF}: paintColor.BgHiRed, // Red Range

		0b101: paintColor.FgHiMagenta, // Magenta M
		//{R: 0x99, G: 0x00, B: 0x99, A: 0xFF}: paintColor.BgHiMagenta, // Magenta M

		0b110: paintColor.FgHiYellow, // Orange
		//		{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}: paintColor.BgHiYellow, //paintColor.FgHiGreen,

		0b111: paintColor.FgHiWhite, // light gray
		//{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}: paintColor.BgHiWhite,

	}

	width       int
	height      int
	cursorReset string
)

type location string

type GUI struct {
	screenBuffer [][]pixel // Buffer to paint stuff in
}

type pixel struct {
	block           string
	colorForeground paintColor.Attribute
	colorBackground paintColor.Attribute
}

type Coords struct {
	X float64
	Y float64
}

func init() {
	width, height, _ = term.GetSize(int(os.Stdin.Fd()))
	cursorReset = fmt.Sprintf("\033[%dD\033[%dA", width, height)
}

// New returns a new GUI Instance
func New() *GUI {
	screenBuffer := make([][]pixel, height)
	for index := range screenBuffer {
		screenBuffer[index] = make([]pixel, width)
	}

	gui := &GUI{
		screenBuffer: screenBuffer,
	}

	gui.FlushScreen()
	return gui
}

// PaintImage paints an image into the screenbuffer
func (gui *GUI) PaintImage(image image.Image, pos Coords, size Coords) {
	screenPos := pos
	var imagePos Coords

	var background paintColor.Attribute
	var foreground paintColor.Attribute

	for imagePos.Y = 0; imagePos.Y < float64(image.Bounds().Max.Y) && screenPos.Y < float64(height); imagePos.Y += float64(image.Bounds().Max.Y) / size.Y {
		for imagePos.X = 0; imagePos.X < float64(image.Bounds().Max.X) && screenPos.X < float64(width); imagePos.X += float64(image.Bounds().Max.X) / size.X {
			if screenPos.X < 0 || screenPos.Y < 0 {
				//screenPos.X += image.Bounds().Max.X / wantedWidth
				screenPos.X++
				continue
			}
			backgroundChan := make(chan paintColor.Attribute)
			foregroundChan := make(chan paintColor.Attribute)

			go gui.GetColor_4bit(location_background, backgroundChan, image, imagePos, screenPos)
			go gui.GetColor_4bit(location_foreground, foregroundChan, image, imagePos, screenPos)

			background = <-backgroundChan
			foreground = <-foregroundChan

			gui.screenBuffer[int(screenPos.Y)][int(screenPos.X)].block = paintColor.New(background, foreground).Sprint("▄")
			gui.screenBuffer[int(screenPos.Y)][int(screenPos.X)].colorBackground = background
			gui.screenBuffer[int(screenPos.Y)][int(screenPos.X)].colorForeground = foreground
			screenPos.X++
		}
		screenPos.Y++
		screenPos.X = pos.X
	}
}

// FlushScreen resets the screen buffer to a white screen
func (gui *GUI) FlushScreen() {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gui.screenBuffer[y][x].block = paintColor.New(paintColor.BgHiWhite, paintColor.FgHiWhite).Sprint("▄")
			gui.screenBuffer[y][x].colorForeground = paintColor.FgHiWhite
			gui.screenBuffer[y][x].colorBackground = paintColor.BgHiWhite
		}
	}
}

// PrintBuffer prints the screenBuffer to terminal
func (gui *GUI) PrintBuffer() {
	var buf bytes.Buffer
	buf.Write([]byte(cursorReset))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			buf.Write([]byte(gui.screenBuffer[y][x].block))
		}
	}

	io.Copy(os.Stdout, &buf)
}

// GetColor returns the best matching color for a pixel
func (gui *GUI) GetColor_4bit(loc location, returnChan chan paintColor.Attribute, image image.Image, imagePos Coords, screenPos Coords) {

	var toBeReturned paintColor.Attribute

	switch loc {
	case location_background:
		r, g, b, a := image.At(int(imagePos.X), int(imagePos.Y)).RGBA()
		if a > 200 {
			shiftR := r >> 15
			shiftG := g >> 15
			shiftB := b >> 15
			shiftedRGB := (shiftR << 2) | (shiftG << 1) | shiftB
			toBeReturned = COLORS_BACKGROUND_4bit[byte(shiftedRGB)]
		} else {
			toBeReturned = gui.screenBuffer[int(screenPos.Y)][int(screenPos.X)].colorBackground
		}

	case location_foreground:
		r, g, b, a := image.At(int(imagePos.X), int(imagePos.Y)+1).RGBA()
		if a > 200 {
			shiftR := r >> 15
			shiftG := g >> 15
			shiftB := b >> 15
			shiftedRGB := (shiftR << 2) | (shiftG << 1) | shiftB
			toBeReturned = COLORS_FOREGROUND_4bit[byte(shiftedRGB)]
		} else {
			toBeReturned = gui.screenBuffer[int(screenPos.Y)][int(screenPos.X)].colorForeground
		}
	}

	returnChan <- toBeReturned
}

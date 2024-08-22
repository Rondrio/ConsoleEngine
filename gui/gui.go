package GUI

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"os"

	paintColor "github.com/fatih/color"
	"golang.org/x/term"
)

const (
	location_foreground = location("foreground")
	location_background = location("background")
)

var (
	COLORS_BACKGROUND = map[color.RGBA]paintColor.Attribute{
		/* A12830 - FF584F*/
		{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}: paintColor.BgRed,   // 	Red Range
		{R: 0x99, G: 0x00, B: 0x00, A: 0xFF}: paintColor.BgHiRed, // Red Range

		{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}: paintColor.BgBlue,
		{R: 0x99, G: 0xCC, B: 0xFF, A: 0xFF}: paintColor.BgHiBlue,

		{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}: paintColor.BgHiYellow, //paintColor.FgHiGreen,
		{R: 0xFF, G: 0x99, B: 0x33, A: 0xFF}: paintColor.BgYellow,   // Orange

		{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}: paintColor.BgBlack,
		{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}: paintColor.BgHiWhite,

		{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}: paintColor.BgCyan,
		{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}: paintColor.BgHiCyan,

		{R: 0x99, G: 0x00, B: 0x4c, A: 0xFF}: paintColor.BgYellow,    // Purple M
		{R: 0x99, G: 0x00, B: 0x99, A: 0xFF}: paintColor.BgHiMagenta, // Magenta M

		{R: 0x29, G: 0x29, B: 0x29, A: 0xFF}: paintColor.BgHiBlack, // Gray M
		{R: 0xE0, G: 0xE0, B: 0xE0, A: 0xFF}: paintColor.BgWhite,   // light gray
	}

	COLORS_FOREGROUND = map[color.RGBA]paintColor.Attribute{
		/* A12830 - FF584F*/
		{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}: paintColor.FgRed,   // 	Red Range
		{R: 0x99, G: 0x00, B: 0x00, A: 0xFF}: paintColor.FgHiRed, // Red Range

		{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}: paintColor.FgBlue,
		{R: 0x99, G: 0xCC, B: 0xFF, A: 0xFF}: paintColor.FgHiBlue,

		{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}: paintColor.FgHiYellow, //paintColor.FgHiGreen,
		{R: 0xFF, G: 0x99, B: 0x33, A: 0xFF}: paintColor.FgYellow,   // Orange

		{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}: paintColor.FgBlack,
		{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}: paintColor.FgHiWhite,

		{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}: paintColor.FgCyan,
		{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}: paintColor.FgHiCyan,

		{R: 0x99, G: 0x00, B: 0x4c, A: 0xFF}: paintColor.FgYellow,    // Purple M
		{R: 0x99, G: 0x00, B: 0x99, A: 0xFF}: paintColor.FgHiMagenta, // Magenta M

		{R: 0x29, G: 0x29, B: 0x29, A: 0xFF}: paintColor.FgHiBlack, // Gray M
		{R: 0xE0, G: 0xE0, B: 0xE0, A: 0xFF}: paintColor.FgWhite,   // light gray
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

// paint image into the screenbuffer
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

			go gui.getColor(location_background, backgroundChan, image, imagePos, screenPos)
			go gui.getColor(location_foreground, foregroundChan, image, imagePos, screenPos)

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

// reset the screen buffer to a white screen
func (gui *GUI) FlushScreen() {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gui.screenBuffer[y][x].block = paintColor.New(paintColor.BgHiWhite, paintColor.FgHiWhite).Sprint("▄")
			gui.screenBuffer[y][x].colorForeground = paintColor.FgHiWhite
			gui.screenBuffer[y][x].colorBackground = paintColor.BgHiWhite
		}
	}
}

// print screenBuffer to terminal
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

// get best matching color for a pixel
func (gui *GUI) getColor(loc location, returnChan chan paintColor.Attribute, image image.Image, imagePos Coords, screenPos Coords) {
	var diff uint32 = 0xFFFFFFFF

	var toBeReturned paintColor.Attribute

	switch loc {
	case location_background:
		r, g, b, a := image.At(int(imagePos.X), int(imagePos.Y)).RGBA()
		if a > 200 {
			for color := range COLORS_BACKGROUND {
				cR, cG, cB, cA := color.RGBA()

				nDiff := uint32(math.Abs(float64(a)-float64(cA)) + math.Abs(float64(r)-float64(cR)) + math.Abs(float64(g)-float64(cG)) + math.Abs(float64(b)-float64(cB)))
				if diff > nDiff {
					diff = nDiff
					toBeReturned = COLORS_BACKGROUND[color]
				}
			}
		} else {
			toBeReturned = gui.screenBuffer[int(screenPos.Y)][int(screenPos.X)].colorBackground
		}

	case location_foreground:
		r, g, b, a := image.At(int(imagePos.X), int(imagePos.Y)+1).RGBA()
		if a > 200 {
			for color := range COLORS_FOREGROUND {
				cR, cG, cB, cA := color.RGBA()

				nDiff := uint32(math.Abs(float64(a)-float64(cA)) + math.Abs(float64(r)-float64(cR)) + math.Abs(float64(g)-float64(cG)) + math.Abs(float64(b)-float64(cB)))
				if diff > nDiff {
					diff = nDiff
					toBeReturned = COLORS_FOREGROUND[color]
				}
			}
		} else {
			toBeReturned = gui.screenBuffer[int(screenPos.Y)][int(screenPos.X)].colorForeground
		}
	}

	returnChan <- toBeReturned
}

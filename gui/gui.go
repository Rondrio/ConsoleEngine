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
	"golang.org/x/crypto/ssh/terminal"
)

const (
	location_foreground = location("foreground")
	location_background = location("background")
)

var (
	COLORS_BACKGROUND = map[color.RGBA]paintColor.Attribute{
		/* A12830 - FF584F*/
		color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}: paintColor.BgRed,   // 	Red Range
		color.RGBA{R: 0x99, G: 0x00, B: 0x00, A: 0xFF}: paintColor.BgHiRed, // Red Range

		color.RGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}: paintColor.BgBlue,
		color.RGBA{R: 0x99, G: 0xCC, B: 0xFF, A: 0xFF}: paintColor.BgHiBlue,

		color.RGBA{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}: paintColor.BgHiYellow, //paintColor.FgHiGreen,
		color.RGBA{R: 0xFF, G: 0x99, B: 0x33, A: 0xFF}: paintColor.BgYellow,   // Orange

		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}: paintColor.BgBlack,
		color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}: paintColor.BgHiWhite,

		color.RGBA{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}: paintColor.BgCyan,
		color.RGBA{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}: paintColor.BgHiCyan,

		color.RGBA{R: 0x99, G: 0x00, B: 0x4c, A: 0xFF}: paintColor.BgYellow,    // Purple M
		color.RGBA{R: 0x99, G: 0x00, B: 0x99, A: 0xFF}: paintColor.BgHiMagenta, // Magenta M

		color.RGBA{R: 0x29, G: 0x29, B: 0x29, A: 0xFF}: paintColor.BgHiBlack, // Gray M
		color.RGBA{R: 0xE0, G: 0xE0, B: 0xE0, A: 0xFF}: paintColor.BgWhite,   // light gray
	}

	COLORS_FOREGROUND = map[color.RGBA]paintColor.Attribute{
		/* A12830 - FF584F*/
		color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}: paintColor.FgRed,   // 	Red Range
		color.RGBA{R: 0x99, G: 0x00, B: 0x00, A: 0xFF}: paintColor.FgHiRed, // Red Range

		color.RGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}: paintColor.FgBlue,
		color.RGBA{R: 0x99, G: 0xCC, B: 0xFF, A: 0xFF}: paintColor.FgHiBlue,

		color.RGBA{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}: paintColor.FgHiYellow, //paintColor.FgHiGreen,
		color.RGBA{R: 0xFF, G: 0x99, B: 0x33, A: 0xFF}: paintColor.FgYellow,   // Orange

		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}: paintColor.FgBlack,
		color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}: paintColor.FgHiWhite,

		color.RGBA{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}: paintColor.FgCyan,
		color.RGBA{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}: paintColor.FgHiCyan,

		color.RGBA{R: 0x99, G: 0x00, B: 0x4c, A: 0xFF}: paintColor.FgYellow,    // Purple M
		color.RGBA{R: 0x99, G: 0x00, B: 0x99, A: 0xFF}: paintColor.FgHiMagenta, // Magenta M

		color.RGBA{R: 0x29, G: 0x29, B: 0x29, A: 0xFF}: paintColor.FgHiBlack, // Gray M
		color.RGBA{R: 0xE0, G: 0xE0, B: 0xE0, A: 0xFF}: paintColor.FgWhite,   // light gray
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
	X int
	Y int
}

func init() {
	width, height, _ = terminal.GetSize(int(os.Stdin.Fd()))
	cursorReset = fmt.Sprintf("\033[%dD\033[%dA", width, height)
}

func New() *GUI {
	screenBuffer := make([][]pixel, height)
	for index := range screenBuffer {
		screenBuffer[index] = make([]pixel, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			screenBuffer[y][x].block = paintColor.New(paintColor.BgHiWhite, paintColor.FgHiWhite).Sprint("▄")
			screenBuffer[y][x].colorForeground = paintColor.FgHiWhite
			screenBuffer[y][x].colorBackground = paintColor.BgHiWhite
		}
	}

	return &GUI{
		screenBuffer: screenBuffer,
	}
}

func (gui *GUI) PaintImage(image image.Image, pos Coords, wantedHeight int, wantedWidth int) {
	screenPos := pos
	var imagePos Coords

	var background paintColor.Attribute
	var foreground paintColor.Attribute

	for imagePos.Y = 0; imagePos.Y < image.Bounds().Dy() && screenPos.Y < height; imagePos.Y += image.Bounds().Dy() / wantedHeight {
		for imagePos.X = 0; imagePos.X < image.Bounds().Dx() && screenPos.X < width; imagePos.X += image.Bounds().Dx() / wantedWidth {

			backgroundChan := make(chan paintColor.Attribute)
			foregroundChan := make(chan paintColor.Attribute)

			go gui.getColor(location_background, backgroundChan, image, imagePos, screenPos)
			go gui.getColor(location_foreground, foregroundChan, image, imagePos, screenPos)

			background = <-backgroundChan
			foreground = <-foregroundChan

			gui.screenBuffer[screenPos.Y][screenPos.X].block = paintColor.New(background, foreground).Sprint("▄")
			gui.screenBuffer[screenPos.Y][screenPos.X].colorBackground = background
			gui.screenBuffer[screenPos.Y][screenPos.X].colorForeground = foreground
			screenPos.X++
		}
		screenPos.Y++
		screenPos.X = pos.X
	}
}

func (g *GUI) PrintBuffer() {
	var buf bytes.Buffer
	buf.Write([]byte(cursorReset))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			buf.Write([]byte(g.screenBuffer[y][x].block))
		}
		buf.Write([]byte("\n"))
	}

	io.Copy(os.Stdout, &buf)
}

func (gui *GUI) getColor(loc location, returnChan chan paintColor.Attribute, image image.Image, imagePos Coords, screenPos Coords) {
	var diff uint32 = 0xFFFFFFFF

	var toBeReturned paintColor.Attribute

	switch loc {
	case location_background:
		r, g, b, a := image.At(imagePos.X, imagePos.Y).RGBA()
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
			toBeReturned = gui.screenBuffer[screenPos.Y][screenPos.X].colorBackground
		}
		returnChan <- toBeReturned

	case location_foreground:
		r, g, b, a := image.At(imagePos.X, imagePos.Y+1).RGBA()
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
			toBeReturned = gui.screenBuffer[screenPos.Y][screenPos.X].colorForeground
		}

		returnChan <- toBeReturned
	}
}

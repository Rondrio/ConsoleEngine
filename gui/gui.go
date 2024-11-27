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
	COLORS_BACKGROUND_4bit_arr = []paintColor.Attribute{
		/*0b000000*/ paintColor.BgBlack,
		/*0b000001*/ paintColor.BgBlue,
		/*0b000010*/ paintColor.BgBlue,
		/*0b000011*/ paintColor.BgHiBlue,
		/*0b000100*/ paintColor.BgGreen,
		/*0b000101*/ paintColor.BgCyan,
		/*0b000110*/ paintColor.BgHiBlue,
		/*0b000111*/ paintColor.BgHiBlue,
		/*0b001000*/ paintColor.BgGreen,
		/*0b001001*/ paintColor.BgGreen,
		/*0b001010*/ paintColor.BgHiCyan,
		/*0b001011*/ paintColor.BgCyan,
		/*0b001100*/ paintColor.BgHiGreen,
		/*0b001101*/ paintColor.BgHiGreen,
		/*0b001110*/ paintColor.BgHiGreen, // maybe cyan?
		/*0b001111*/ paintColor.BgHiCyan,
		/*0b010000*/ paintColor.BgRed,
		/*0b010001*/ paintColor.BgHiMagenta,
		/*0b010010*/ paintColor.BgMagenta,
		/*0b010011*/ paintColor.BgBlue, // maybe magenta
		/*0b010100*/ paintColor.BgYellow,
		/*0b010101*/ paintColor.BgHiBlack, // dark gray
		/*0b010110*/ paintColor.BgHiBlue, // dark gray
		/*0b010111*/ paintColor.BgHiBlue,
		/*0b011000*/ paintColor.BgHiGreen,
		/*0b011001*/ paintColor.BgGreen,
		/*0b011010*/ paintColor.BgCyan,
		/*0b011011*/ paintColor.BgHiBlue,
		/*0b011100*/ paintColor.BgGreen,
		/*0b011101*/ paintColor.BgGreen,
		/*0b011110*/ paintColor.BgHiGreen,
		/*0b011111*/ paintColor.BgHiCyan,
		/*0b100000*/ paintColor.BgRed,
		/*0b100001*/ paintColor.BgMagenta,
		/*0b100010*/ paintColor.BgMagenta,
		/*0b100011*/ paintColor.BgMagenta,
		/*0b100100*/ paintColor.BgYellow,
		/*0b100101*/ paintColor.BgRed, // maybe brown/yellow
		/*0b100110*/ paintColor.BgHiMagenta,
		/*0b100111*/ paintColor.BgHiMagenta,
		/*0b101000*/ paintColor.BgHiYellow,
		/*0b101001*/ paintColor.BgHiYellow,
		/*0b101010*/ paintColor.BgWhite,
		/*0b101011*/ paintColor.BgHiCyan,
		/*0b101100*/ paintColor.BgGreen,
		/*0b101101*/ paintColor.BgGreen,
		/*0b101110*/ paintColor.BgGreen,
		/*0b101111*/ paintColor.BgCyan,
		/*0b110000*/ paintColor.BgHiRed, // 	Red Range
		/*0b110001*/ paintColor.BgHiRed,
		/*0b110010*/ paintColor.BgHiMagenta,
		/*0b110011*/ paintColor.BgHiMagenta, // Magenta M
		/*0b110100*/ paintColor.BgHiRed, // 	Red Range
		/*0b110101*/ paintColor.BgHiRed, // 	Red Range
		/*0b110110*/ paintColor.BgHiMagenta, // 	Red Range
		/*0b110111*/ paintColor.BgHiMagenta, // 	Red Range
		/*0b111000*/ paintColor.BgYellow, // 	Red Range
		/*0b111001*/ paintColor.BgYellow, // 	Red Range
		/*0b111010*/ paintColor.BgHiRed, // 	Red Range
		/*0b111011*/ paintColor.BgHiMagenta, // 	Red Range
		/*0b111100*/ paintColor.BgHiYellow, // Orange
		/*0b111101*/ paintColor.BgHiYellow, // Orange
		/*0b111110*/ paintColor.BgHiYellow, // Orange
		/*0b111111*/ paintColor.BgHiWhite, // light gray
	}

	COLORS_FOREGROUND_4bit_arr = []paintColor.Attribute{
		/*0b000000*/ paintColor.FgBlack,
		/*0b000001*/ paintColor.FgBlue,
		/*0b000010*/ paintColor.FgBlue,
		/*0b000011*/ paintColor.FgHiBlue,
		/*0b000100*/ paintColor.FgGreen,
		/*0b000101*/ paintColor.FgCyan,
		/*0b000110*/ paintColor.FgHiBlue,
		/*0b000111*/ paintColor.FgHiBlue,
		/*0b001000*/ paintColor.FgGreen,
		/*0b001001*/ paintColor.FgGreen,
		/*0b001010*/ paintColor.FgHiCyan,
		/*0b001011*/ paintColor.FgCyan,
		/*0b001100*/ paintColor.FgHiGreen,
		/*0b001101*/ paintColor.FgHiGreen,
		/*0b001110*/ paintColor.FgHiGreen, // maybe cyan?
		/*0b001111*/ paintColor.FgHiCyan,
		/*0b010000*/ paintColor.FgRed,
		/*0b010001*/ paintColor.FgHiMagenta,
		/*0b010010*/ paintColor.FgMagenta,
		/*0b010011*/ paintColor.FgBlue, // maybe magenta
		/*0b010100*/ paintColor.FgYellow,
		/*0b010101*/ paintColor.FgHiBlack, // dark gray
		/*0b010110*/ paintColor.FgHiBlue, // dark gray
		/*0b010111*/ paintColor.FgHiBlue,
		/*0b011000*/ paintColor.FgHiGreen,
		/*0b011001*/ paintColor.FgGreen,
		/*0b011010*/ paintColor.FgCyan,
		/*0b011011*/ paintColor.FgHiBlue,
		/*0b011100*/ paintColor.FgGreen,
		/*0b011101*/ paintColor.FgGreen,
		/*0b011110*/ paintColor.FgHiGreen,
		/*0b011111*/ paintColor.FgHiCyan,
		/*0b100000*/ paintColor.FgRed,
		/*0b100001*/ paintColor.FgMagenta,
		/*0b100010*/ paintColor.FgMagenta,
		/*0b100011*/ paintColor.FgMagenta,
		/*0b100100*/ paintColor.FgYellow,
		/*0b100101*/ paintColor.FgRed, // maybe brown/yellow
		/*0b100110*/ paintColor.FgHiMagenta,
		/*0b100111*/ paintColor.FgHiMagenta,
		/*0b101000*/ paintColor.FgHiYellow,
		/*0b101001*/ paintColor.FgHiYellow,
		/*0b101010*/ paintColor.FgWhite,
		/*0b101011*/ paintColor.FgHiCyan,
		/*0b101100*/ paintColor.FgGreen,
		/*0b101101*/ paintColor.FgGreen,
		/*0b101110*/ paintColor.FgGreen,
		/*0b101111*/ paintColor.FgCyan,
		/*0b110000*/ paintColor.FgHiRed, // 	Red Range
		/*0b110001*/ paintColor.FgHiRed,
		/*0b110010*/ paintColor.FgHiMagenta,
		/*0b110011*/ paintColor.FgHiMagenta, // Magenta M
		/*0b110100*/ paintColor.FgHiRed, // 	Red Range
		/*0b110101*/ paintColor.FgHiRed, // 	Red Range
		/*0b110110*/ paintColor.FgHiMagenta, // 	Red Range
		/*0b110111*/ paintColor.FgHiMagenta, // 	Red Range
		/*0b111000*/ paintColor.FgYellow, // 	Red Range
		/*0b111001*/ paintColor.FgYellow, // 	Red Range
		/*0b111010*/ paintColor.FgHiRed, // 	Red Range
		/*0b111011*/ paintColor.FgHiMagenta, // 	Red Range
		/*0b111100*/ paintColor.FgHiYellow, // Orange
		/*0b111101*/ paintColor.FgHiYellow, // Orange
		/*0b111110*/ paintColor.FgHiYellow, // Orange
		/*0b111111*/ paintColor.FgHiWhite, // light gray
	}

	COLORS_BACKGROUND_4bit = map[byte]paintColor.Attribute{
		// black
		0b000000: paintColor.BgBlack,
		0b010101: paintColor.BgHiBlack, // dark gray

		// blue
		0b000001: paintColor.BgBlue,
		0b000010: paintColor.BgBlue,
		0b000011: paintColor.BgHiBlue,
		0b000111: paintColor.BgHiBlue,
		0b010011: paintColor.BgBlue, // maybe magenta
		0b010111: paintColor.BgHiBlue,
		0b011011: paintColor.BgHiBlue,

		// green
		0b000100: paintColor.BgGreen,
		0b001000: paintColor.BgGreen,
		0b001001: paintColor.BgGreen,
		0b001100: paintColor.BgHiGreen,
		0b001101: paintColor.BgHiGreen,
		0b001110: paintColor.BgHiGreen, // maybe cyan?
		0b011000: paintColor.BgHiGreen,
		0b011001: paintColor.BgGreen,
		0b011100: paintColor.BgGreen,
		0b101100: paintColor.BgGreen,
		0b101101: paintColor.BgGreen,
		0b101110: paintColor.BgGreen,
		0b011101: paintColor.BgGreen,

		// cyan
		0b000101: paintColor.BgCyan,
		0b001010: paintColor.BgHiCyan,
		0b001011: paintColor.BgCyan,
		0b001111: paintColor.BgHiCyan,
		0b011010: paintColor.BgCyan,
		0b101111: paintColor.BgCyan,
		0b101011: paintColor.BgHiCyan,
		0b011111: paintColor.BgHiCyan,

		// red
		0b010000: paintColor.BgRed,
		0b100101: paintColor.BgRed,   // maybe brown/yellow
		0b110000: paintColor.BgHiRed, // 	Red Range
		0b110001: paintColor.BgHiRed,
		0b110100: paintColor.BgHiRed, // 	Red Range
		0b110101: paintColor.BgHiRed, // 	Red Range
		0b111010: paintColor.BgHiRed, // 	Red Range
		0b100000: paintColor.BgRed,

		// magenta
		0b010001: paintColor.BgHiMagenta,
		0b010010: paintColor.BgMagenta,
		0b100001: paintColor.BgMagenta,
		0b100010: paintColor.BgMagenta,
		0b100011: paintColor.BgMagenta,
		0b100110: paintColor.BgHiMagenta,
		0b110010: paintColor.BgHiMagenta,
		0b100111: paintColor.BgHiMagenta,
		0b110110: paintColor.BgHiMagenta, // 	Red Range
		0b110111: paintColor.BgHiMagenta, // 	Red Range
		0b111011: paintColor.BgHiMagenta, // 	Red Range
		0b110011: paintColor.BgHiMagenta, // Magenta M

		// yellow
		0b010100: paintColor.BgYellow,
		0b100100: paintColor.BgYellow,
		0b101000: paintColor.BgHiYellow,
		0b101001: paintColor.BgHiYellow,
		0b111000: paintColor.BgYellow,   // 	Red Range
		0b111001: paintColor.BgYellow,   // 	Red Range
		0b111101: paintColor.BgHiYellow, // Orange
		0b111110: paintColor.BgHiYellow, // Orange
		0b111100: paintColor.BgHiYellow, // Orange

		// white
		0b101010: paintColor.BgWhite,
		0b111111: paintColor.BgHiWhite, // light gray
	}

	COLORS_FOREGROUND_4bit = map[byte]paintColor.Attribute{
		0b000000: paintColor.FgBlack,
		0b010101: paintColor.FgHiBlack, // dark gray

		// blue
		0b000001: paintColor.FgBlue,
		0b000010: paintColor.FgBlue,
		0b000011: paintColor.FgHiBlue,
		0b000111: paintColor.FgHiBlue,
		0b010011: paintColor.FgBlue, // maybe magenta
		0b010111: paintColor.FgHiBlue,
		0b011011: paintColor.FgHiBlue,

		// green
		0b000100: paintColor.FgGreen,
		0b001000: paintColor.FgGreen,
		0b001001: paintColor.FgGreen,
		0b001100: paintColor.FgHiGreen,
		0b001101: paintColor.FgHiGreen,
		0b001110: paintColor.FgHiGreen, // maybe cyan?
		0b011000: paintColor.FgHiGreen,
		0b011001: paintColor.FgGreen,
		0b011100: paintColor.FgGreen,
		0b101100: paintColor.FgGreen,
		0b101101: paintColor.FgGreen,
		0b101110: paintColor.FgGreen,
		0b011101: paintColor.FgGreen,

		// cyan
		0b000101: paintColor.FgCyan,
		0b001010: paintColor.FgHiCyan,
		0b001011: paintColor.FgCyan,
		0b001111: paintColor.FgHiCyan,
		0b011010: paintColor.FgCyan,
		0b101111: paintColor.FgCyan,
		0b101011: paintColor.FgHiCyan,
		0b011111: paintColor.FgHiCyan,

		// red
		0b010000: paintColor.FgRed,
		0b100101: paintColor.FgRed,   // maybe brown/yellow
		0b110000: paintColor.FgHiRed, // 	Red Range
		0b110001: paintColor.FgHiRed,
		0b110100: paintColor.FgHiRed, // 	Red Range
		0b110101: paintColor.FgHiRed, // 	Red Range
		0b111010: paintColor.FgHiRed, // 	Red Range
		0b100000: paintColor.FgHiRed,

		// magenta
		0b010001: paintColor.FgHiMagenta,
		0b010010: paintColor.FgMagenta,
		0b100001: paintColor.FgMagenta,
		0b100010: paintColor.FgMagenta,
		0b100011: paintColor.FgMagenta,
		0b100110: paintColor.FgHiMagenta,
		0b110010: paintColor.FgHiMagenta,
		0b100111: paintColor.FgHiMagenta,
		0b110110: paintColor.FgHiMagenta, // 	Red Range
		0b110111: paintColor.FgHiMagenta, // 	Red Range
		0b111011: paintColor.FgHiMagenta, // 	Red Range
		0b110011: paintColor.FgHiMagenta, // Magenta M

		// yellow
		0b010100: paintColor.FgYellow,
		0b100100: paintColor.FgYellow,
		0b101000: paintColor.FgHiYellow,
		0b101001: paintColor.FgHiYellow,
		0b111000: paintColor.FgYellow,   // 	Red Range
		0b111001: paintColor.FgYellow,   // 	Red Range
		0b111101: paintColor.FgHiYellow, // Orange
		0b111110: paintColor.FgHiYellow, // Orange
		0b111100: paintColor.FgHiYellow, // Orange

		// white
		0b101010: paintColor.FgWhite,
		0b111111: paintColor.FgHiWhite, // light gray
	}

	width       int
	height      int
	cursorReset string
	whiteBlock  = paintColor.New(paintColor.BgHiWhite, paintColor.FgHiWhite).Sprint("▄")
	buf         bytes.Buffer
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

func (gui *GUI) GetScreenSize() Coords {
	return Coords{float64(width), float64(height)}
}

// PaintImage paints an image into the screenbuffer
func (gui *GUI) PaintImage(image image.Image, pos Coords, size Coords) {
	screenPos := pos
	var imagePos Coords

	for imagePos.Y = 0; imagePos.Y < float64(image.Bounds().Max.Y) && screenPos.Y < float64(height); imagePos.Y += float64(image.Bounds().Max.Y) / size.Y {
		for imagePos.X = 0; imagePos.X < float64(image.Bounds().Max.X) && screenPos.X < float64(width); imagePos.X += float64(image.Bounds().Max.X) / size.X {
			if screenPos.X < 0 || screenPos.Y < 0 {
				//screenPos.X += image.Bounds().Max.X / wantedWidth
				screenPos.X++
				continue
			}
			paintPixel(gui, image, imagePos, screenPos)
			screenPos.X++
		}
		screenPos.Y++
		screenPos.X = pos.X
	}
}

func paintPixel(gui *GUI, image image.Image, imagePos Coords, screenPos Coords) {
	background := gui.GetColor_4bit(location_background, image, imagePos, screenPos)
	foreground := gui.GetColor_4bit(location_foreground, image, imagePos, screenPos)

	gui.screenBuffer[int(screenPos.Y)][int(screenPos.X)].block = paintColor.New(background, foreground).Sprint("▄")
	gui.screenBuffer[int(screenPos.Y)][int(screenPos.X)].colorBackground = background
	gui.screenBuffer[int(screenPos.Y)][int(screenPos.X)].colorForeground = foreground
}

// FlushScreen resets the screen buffer to a white screen
func (gui *GUI) FlushScreen() {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gui.screenBuffer[y][x].block = whiteBlock
			gui.screenBuffer[y][x].colorForeground = paintColor.FgHiWhite
			gui.screenBuffer[y][x].colorBackground = paintColor.BgHiWhite
		}
	}
}

// PrintBuffer prints the screenBuffer to terminal
func (gui *GUI) PrintBuffer() {
	buf.Write([]byte(cursorReset))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			buf.Write([]byte(gui.screenBuffer[y][x].block))
		}
	}

	io.Copy(os.Stdout, &buf)
	buf.Truncate(0)
}

// GetColor returns the best matching color for a pixel
func (gui *GUI) GetColor_4bit(loc location, image image.Image, imagePos Coords, screenPos Coords) paintColor.Attribute {

	var toBeReturned paintColor.Attribute

	switch loc {
	case location_background:
		r, g, b, a := image.At(int(imagePos.X), int(imagePos.Y)).RGBA()
		if a > 0 {
			shiftR := r >> 14
			shiftG := g >> 14
			shiftB := b >> 14
			shiftedRGB := (shiftR << 4) | (shiftG << 2) | shiftB
			toBeReturned = COLORS_BACKGROUND_4bit_arr[byte(shiftedRGB)]
		} else {
			toBeReturned = gui.screenBuffer[int(screenPos.Y)][int(screenPos.X)].colorBackground
		}

	case location_foreground:
		r, g, b, a := image.At(int(imagePos.X), int(imagePos.Y)+1).RGBA()
		if a > 0 {
			shiftR := r >> 14
			shiftG := g >> 14
			shiftB := b >> 14
			shiftedRGB := (shiftR << 4) | (shiftG << 2) | shiftB
			toBeReturned = COLORS_FOREGROUND_4bit_arr[byte(shiftedRGB)]
		} else {
			toBeReturned = gui.screenBuffer[int(screenPos.Y)][int(screenPos.X)].colorForeground
		}
	}

	return toBeReturned
}

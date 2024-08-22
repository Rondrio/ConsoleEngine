package GUI

import (
	"image/color"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"golang.org/x/term"
)

type Screen struct {
	GUI          *GUI
	Layers       map[int]map[string]Entity
	ScreenSize   Coords
	KeyEventChan chan KeyCode
}

func NewScreen(gui *GUI) *Screen {

	screen := &Screen{
		GUI:          gui,
		Layers:       make(map[int]map[string]Entity),
		ScreenSize:   Coords{float64(width), float64(height)},
		KeyEventChan: make(chan KeyCode),
	}

	screen.Layers = make(map[int]map[string]Entity)

	go screen.readInput()

	return screen
}

func (screen *Screen) Paint() {

	max := max(screen.Layers)

	screen.GUI.FlushScreen()

	for i := 0; i <= max; i++ {
		if len(screen.Layers[i]) != 0 {
			for _, entity := range screen.Layers[i] {
				if entity.GetVisible() {
					bg := color.RGBA{0, 0, 0, 0}
					texture := imaging.Rotate(entity.GetTexture(), float64(entity.GetRotation()), bg)
					screen.GUI.PaintImage(texture, entity.GetPosition(), entity.GetSize())
				}
			}
		}
	}

	screen.GUI.PrintBuffer()
}

func (screen *Screen) AddEntity(entity Entity) {
	if _, ok := screen.Layers[entity.GetLayer()]; !ok {
		screen.Layers[entity.GetLayer()] = make(map[string]Entity)
	}
	screen.Layers[entity.GetLayer()][entity.GetId()] = entity
}

func (screen *Screen) RemoveEntity(toBeDeleted Entity) {
	for layerId := range screen.Layers {
		delete(screen.Layers[layerId], toBeDeleted.GetId())
	}
}

func (screen *Screen) FindEntity(id string) Entity {
	for _, layer := range screen.Layers {
		if entity, ok := layer[id]; ok {
			return entity
		}
	}
	return nil
}

func max(entities map[int]map[string]Entity) int {
	var maxNumber int
	for maxNumber = range entities {
		break
	}
	for n := range entities {
		if n > maxNumber {
			maxNumber = n
		}
	}
	return maxNumber
}

func (screen *Screen) readInput() {

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Println(err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	var b []byte = make([]byte, 1)

	for {
		os.Stdin.Read(b)
		switch KeyCode(b[0]) {
		case escape:
			term.Restore(int(os.Stdin.Fd()), oldState)
			os.Exit(1)
		default:
			screen.KeyEventChan <- KeyCode(b[0])
		}
	}
}

package Engine2D

import (
	GUI "consoleEngine/gui"
	"image/color"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"golang.org/x/term"
)

type Screen struct {
	GUI          *GUI.GUI
	Layers       map[int]map[string]Entity
	ScreenSize   GUI.Coords
	KeyEventChan chan KeyCode
}

// NewScreen() creates and returns a new Screen on the GUI
func NewScreen(gui *GUI.GUI) *Screen {
	width, height, _ := term.GetSize(int(os.Stdin.Fd()))

	screen := &Screen{
		GUI:          gui,
		Layers:       make(map[int]map[string]Entity),
		ScreenSize:   GUI.Coords{float64(width), float64(height)},
		KeyEventChan: make(chan KeyCode),
	}

	screen.Layers = make(map[int]map[string]Entity)

	go screen.readInput()

	return screen
}

// Paint() paints the screen on the GUI
func (screen *Screen) Paint() {

	max := getHightestLayer(screen.Layers)

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

// AddEntity() adds a entity to the screen
func (screen *Screen) AddEntity(entity Entity) {
	if _, ok := screen.Layers[entity.GetLayer()]; !ok {
		screen.Layers[entity.GetLayer()] = make(map[string]Entity)
	}
	screen.Layers[entity.GetLayer()][entity.GetId()] = entity
}

// RemoveEntity() removes an entity from the screen
func (screen *Screen) RemoveEntity(toBeDeleted Entity) {
	for layerId := range screen.Layers {
		delete(screen.Layers[layerId], toBeDeleted.GetId())
	}
}

// FindEntity returns the entity with a given string
func (screen *Screen) FindEntity(id string) Entity {
	for _, layer := range screen.Layers {
		if entity, ok := layer[id]; ok {
			return entity
		}
	}
	return nil
}

// getHightestLayer returns the highest layer
func getHightestLayer(layers map[int]map[string]Entity) int {
	var highestLayer int
	// set highestLayer to the first layer in map
	for highestLayer = range layers {
		break
	}

	// iterate through layers, set highest layer to the hightest number
	for n := range layers {
		if n > highestLayer {
			highestLayer = n
		}
	}
	return highestLayer
}

// readInput listens for input events in the console and pushes them into the KeyEventChan
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

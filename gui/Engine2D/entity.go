package Engine2D

import (
	GUI "consoleEngine/gui"
	"image"

	"github.com/google/uuid"
)

type EntityInputHandler func(Entity, KeyCode)

type Entity interface {
	GetId() string
	GetScreen() *Screen

	SetTexture(image.Image)
	GetTexture() image.Image

	SetPosition(pos GUI.Coords)
	GetPosition() GUI.Coords

	SetSize(size GUI.Coords)
	GetSize() GUI.Coords

	SetLayer(layer int)
	GetLayer() int

	SetVisible(visible bool)
	GetVisible() bool

	SetRotation(rotation float64)
	GetRotation() float64
}

type BaseEntity struct {
	Screen   *Screen
	Id       string
	Texture  image.Image
	Pos      GUI.Coords
	Size     GUI.Coords
	Layer    int
	Visible  bool
	Rotation float64
}

func NewEntity(screen *Screen, texture image.Image, pos GUI.Coords, size GUI.Coords, layer int, inputHandler EntityInputHandler) *BaseEntity {
	sprite := BaseEntity{
		Screen:   screen,
		Id:       uuid.NewString(),
		Texture:  texture,
		Pos:      pos,
		Size:     size,
		Layer:    layer,
		Visible:  true,
		Rotation: 0,
	}
	if inputHandler != nil {
		go sprite.readInput(screen, inputHandler)
	}
	return &sprite
}

func (entity *BaseEntity) readInput(screen *Screen, inputHandler EntityInputHandler) {
	for {
		key := <-screen.KeyEventChan
		inputHandler(entity, key)
	}
}

func (entity *BaseEntity) GetId() string {
	return entity.Id
}

func (entity *BaseEntity) GetScreen() *Screen {
	return entity.Screen
}

func (entity *BaseEntity) SetTexture(texture image.Image) {
	entity.Texture = texture
}

func (entity *BaseEntity) GetTexture() image.Image {
	return entity.Texture
}

func (entity *BaseEntity) SetPosition(pos GUI.Coords) {
	entity.Pos = pos
}

func (entity *BaseEntity) GetPosition() GUI.Coords {
	return entity.Pos
}

func (entity *BaseEntity) SetSize(size GUI.Coords) {
	entity.Size = size
}

func (entity *BaseEntity) GetSize() GUI.Coords {
	return entity.Size
}

func (entity *BaseEntity) SetLayer(layer int) {
	entity.Screen.RemoveEntity(entity)
	entity.Layer = layer
	entity.Screen.AddEntity(entity)
}

func (entity *BaseEntity) GetLayer() int {
	return entity.Layer
}

func (entity *BaseEntity) SetVisible(visible bool) {
	entity.Visible = visible
}

func (entity *BaseEntity) GetVisible() bool {
	return entity.Visible
}

func (entity *BaseEntity) SetRotation(rotation float64) {
	entity.Rotation = rotation
}

func (entity *BaseEntity) GetRotation() float64 {
	return entity.Rotation
}

package Engine2D

import (
	GUI "consoleEngine/gui"
	"image"
	"time"
)

const (
	gravitation  = 0.5
	physics_tick = 1000 / 60 * time.Millisecond // time between physics calculations in ms
)

type PhysicsEntityInputHandler func(IPhysicsEntity, KeyCode)

type IPhysicsEntity interface {
	Entity

	SetCollision(collision bool)
	GetCollision() bool

	SetMomentum(momentum GUI.Coords)
	GetMomentum() GUI.Coords

	SetWeight(weight float64)
	GetWeight() float64
}

type PhysicsEntity struct {
	*BaseEntity
	Collision      bool
	Momentum       GUI.Coords
	Weight         float64
	PhysicsEnabled bool
	Grounded       bool
	InputHandler   PhysicsEntityInputHandler
}

func NewPhysicsEntity(screen *Screen, texture image.Image, pos GUI.Coords, size GUI.Coords, layer int, inputHandler PhysicsEntityInputHandler, collision bool, momentum GUI.Coords, weight float64) *PhysicsEntity {
	entity := NewEntity(screen, texture, pos, size, layer, nil)
	physicsEntity := PhysicsEntity{
		BaseEntity:     entity,
		Collision:      collision,
		Momentum:       momentum,
		Weight:         weight,
		PhysicsEnabled: true,
		Grounded:       true,
	}

	if inputHandler != nil {
		go physicsEntity.readInput(screen, inputHandler)
	}

	go physicsEntity.ApplyPhysics()

	return &physicsEntity
}

func (entity *PhysicsEntity) ApplyPhysics() {
	screen := entity.Screen
	ticker := time.NewTicker(physics_tick)

	for {
		if entity.PhysicsEnabled {

			pos := entity.GetPosition()
			mom := entity.GetMomentum()
			// accelerate gravity
			if pos.Y+entity.Size.Y+mom.Y < screen.ScreenSize.Y {
				mom.Y += gravitation
				entity.Grounded = false
			}

			if !entity.Grounded && pos.Y+entity.Size.Y+mom.Y > screen.ScreenSize.Y {
				entity.Grounded = true
				mom.Y = 0
			}

			if mom.X > 0 && entity.Grounded {
				mom.X -= 1
			} else if mom.X < 0 && entity.Grounded {
				mom.X += 1
			}

			// bounce from sides
			if pos.X+entity.Size.X+mom.X >= screen.ScreenSize.X {
				mom.X = -mom.X
			} else if pos.X+mom.X < 0 {
				mom.X = -mom.X
			}

			/*	May be stupid IDK
				// bounce from bottom
				if entity.Pos.Y+entity.Size.Y >= screen.ScreenSize.Y && entity.Momentum.Y > 0 {
					entity.Momentum.Y = -(1 - entity.Weight/100) * entity.Momentum.Y
				}
			*/
			pos.Y += mom.Y
			pos.X += mom.X

			entity.SetMomentum(mom)
			entity.SetPosition(pos)
		}
		<-ticker.C
	}
}

func (entity *PhysicsEntity) readInput(screen *Screen, inputHandler PhysicsEntityInputHandler) {
	for {
		key := <-screen.KeyEventChan
		inputHandler(entity, key)
	}
}

func (pe *PhysicsEntity) SetCollision(collision bool) {
	pe.Collision = collision
}

func (pe *PhysicsEntity) GetCollision() bool {
	return pe.Collision
}

func (pe *PhysicsEntity) SetMomentum(momentum GUI.Coords) {
	pe.Momentum = momentum
}

func (pe *PhysicsEntity) GetMomentum() GUI.Coords {
	return pe.Momentum
}

func (pe *PhysicsEntity) SetWeight(weight float64) {
	pe.Weight = weight
}

func (pe *PhysicsEntity) GetWeight() float64 {
	return pe.Weight
}

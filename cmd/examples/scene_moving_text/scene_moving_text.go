package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/gravestench/akara"
	director "github.com/gravestench/director/pkg"
)

const key = "Director Example - Moving Text"

const numTextObjects = 100

type MovingTextScene struct {
	director.Scene
	textObjects [numTextObjects]akara.EID
	Velocity VelocityFactory
}

func (scene *MovingTextScene) Key() string {
	return key
}

func (scene *MovingTextScene) IsInitialized() bool {
	if scene.Director.World == nil {
		return false
	}

	return true
}

func (scene *MovingTextScene) Init(w *akara.World) {
	fmt.Println("moving text scene init")

	rand.Seed(time.Now().UnixNano())

	scene.makeLabels()

	scene.Director.Window.Title = scene.Key()
	scene.Director.World = w
	scene.BaseSystem.World = w
	scene.InjectComponent(&Velocity{}, &scene.Velocity.ComponentFactory)
}

func (scene *MovingTextScene) makeLabels() {
	ww, wh := scene.Window.Width, scene.Window.Height

	fontSize := wh / 10
	fontName := "Sans"

	for idx := range scene.textObjects {
		rx, ry := rand.Intn(ww), rand.Intn(wh)
		scene.textObjects[idx] = scene.Add.Label("", rx, ry, fontSize, fontName, randColor())
	}
}

func (scene *MovingTextScene) Update(delta time.Duration) {
	scene.updateString()

	for _, eid := range scene.textObjects {
		scene.updatePosition(eid)
	}
}

func (scene *MovingTextScene) updateString() {
	for _, e := range scene.textObjects {
		text, found := scene.Text.Get(e)
		if !found {
			continue
		}

		uuid, found := scene.UUID.Get(e)
		if !found {
			continue
		}

		text.String = uuid.String()[:4]
	}
}

func (scene *MovingTextScene) updatePosition(eid akara.EID) {
	trs, found := scene.Transform.Get(eid)
	if !found  {
		return
	}

	position := trs.Translation

	velocity, found := scene.Velocity.Get(eid)
	if !found  {
		velocity = scene.Velocity.Add(eid)
	}

	const max = 8

	position.X += float64(velocity.X)
	position.Y += float64(velocity.Y)

	velocity.X += (rand.Float32() * 2) - 1
	velocity.Y += (rand.Float32() * 2) - 1

	velocity.X = float32(clamp(velocity.X, -max, max))
	velocity.Y = float32(clamp(velocity.Y, -max, max))

	ww, wh := float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
	tw, th := float32(400), float32(140)

	position.X = float64(wrap(float32(position.X), -tw, ww))
	position.Y = float64(wrap(float32(position.Y), -th, wh))
}

func clamp(v, min, max float32) float32 {
	if v > max {
		v = max
	} else if v < min {
		v = min
	}

	return v
}

func wrap(v, min, max float32) float32 {
	if v > max {
		v = min
	} else if v < min {
		v = max
	}

	return v
}

func randColor() color.Color {
	return &color.RGBA{
		R: uint8(rand.Intn(math.MaxUint8)),
		G: uint8(rand.Intn(math.MaxUint8)),
		B: uint8(rand.Intn(math.MaxUint8)),
		A: uint8(rand.Intn(math.MaxUint8)),
	}
}


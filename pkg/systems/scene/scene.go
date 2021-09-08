package scene

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gravestench/akara"
	director "github.com/gravestench/director/pkg"
	"github.com/gravestench/director/pkg/common"
	"github.com/gravestench/director/pkg/components"
	"github.com/gravestench/mathlib"
	"github.com/gravestench/scenegraph"
	lua "github.com/yuin/gopher-lua"
	"time"
)

type Scene struct {
	akara.BaseSystem
	*director.Director
	Lua         *lua.LState
	Components  common.BasicComponents
	Graph       scenegraph.Node
	key         string
	Add         ObjectFactory
	Renderables *akara.Subscription
	Cameras     []common.Entity
	Width       int
	Height      int
}

var tmpVect mathlib.Vector3

func (s *Scene) renderEntity(e common.Entity) {
	texture, textureFound := s.Components.Texture2D.Get(e)
	rt, rtFound := s.Components.RenderTexture2D.Get(e)
	if !textureFound && !rtFound {
		return
	}

	var t *rl.Texture2D

	if !rtFound {
		t = &texture.Texture2D
	} else {
		t = &rt.Texture
	}

	trs, found := s.Components.Transform.Get(e)
	if !found {
		return
	}

	origin, found := s.Components.Origin.Get(e)
	if !found {
		return
	}

	tmpVect.Set(float64(t.Width), float64(t.Height), 1)

	yRad := trs.Rotation.Y * mathlib.DegreesToRadians
	ov2 := mathlib.NewVector2(origin.Clone().Multiply(&tmpVect).XY()).Rotate(yRad).Negate()
	ov3 := mathlib.NewVector3(ov2.X, ov2.Y, 0)

	x, y := trs.Translation.Clone().Add(ov3).XY()
	v2 := mathlib.NewVector2(x, y)

	position := rl.Vector2{
		X: float32(v2.X),
		Y: float32(v2.Y),
	}

	rotation := float32(trs.Rotation.Y)

	scale := float32(trs.Scale.X)

	rl.DrawTextureEx(*t, position, rotation, scale, rl.White)
}

func (s *Scene) Initialize(d *director.Director, width, height int) {
	s.Add.scene = s
	s.Width = width
	s.Height = height
	s.Director = d
	s.Components.Init(s.Director.World)

	filter := s.Director.World.NewComponentFilter()
	filter.Require(&components.Transform{})
	filter.RequireOne(&components.RenderTexture2D{}, &components.Texture2D{})
	filter.Forbid(&components.Viewport{})

	s.Renderables = s.Director.AddSubscription(filter.Build())
}

func (s *Scene) InitializeLua() {
	if s.LuaInitialized() {
		return
	}

	s.Lua = lua.NewState()

	var luaTypeExporters = []func(*Scene) common.LuaTypeExport{
		luaRectangleTypeExporter,
		luaCircleTypeExporter,
		luaImageTypeExporter,
	}

	for _, luaTypeExporter := range luaTypeExporters {
		luaTypeExport := luaTypeExporter(s)
		common.RegisterLuaType(s.Lua, luaTypeExport)
	}

	s.initComponentsTable()
}

func (s *Scene) initComponentsTable() {
	componentsTable := s.Lua.NewTable()
	s.Lua.SetGlobal("components", componentsTable)

	s.luaExportComponentTransform(componentsTable)
	s.luaExportComponentOrigin(componentsTable)
}

func (s *Scene) UninitializeLua() {
	s.Lua = nil
}

func (s *Scene) LuaInitialized() bool {
	return s.Lua != nil
}

func (s *Scene) updateSceneGraph() {
	for _, eid := range s.Renderables.GetEntities() {
		node, found := s.Components.SceneGraphNode.Get(eid)
		if !found {
			continue
		}

		trs, found := s.Components.Transform.Get(eid)
		if !found {
			continue
		}

		node.Local = trs.GetMatrix()
	}

	s.Graph.UpdateWorldMatrix()
}

func (s *Scene) updateSceneObjects(dt time.Duration) {
	s.Add.update(dt)
}

func (s *Scene) GenericUpdate(dt time.Duration) {
	s.updateSceneGraph()
	s.updateSceneObjects(dt)
}

func (s *Scene) Render() {
	if len(s.Cameras) < 1 {
		s.initCamera()
	}

	for _, cameraEntity := range s.Cameras {
		s.renderEntitiesWithRespectToCamera(cameraEntity)
	}
}

func (s *Scene) initCamera() {
	s.Cameras = make([]common.Entity, 0)
	s.Cameras = append(s.Cameras, s.Add.Camera(0, 0, s.Width, s.Height))
}

func (s *Scene) renderEntitiesWithRespectToCamera(camera common.Entity) {
	rt, found := s.Components.RenderTexture2D.Get(camera)
	if !found {
		return
	}

	cam, found := s.Components.Viewport.Get(camera)
	if !found {
		return
	}

	rl.BeginTextureMode(*rt.RenderTexture2D)
	r, g, b, a := cam.Background.RGBA()
	rl.ClearBackground(rl.NewColor(uint8(r), uint8(g), uint8(b), uint8(a)))
	rl.BeginMode2D(cam.Camera2D)

	for _, entity := range s.Renderables.GetEntities() {
		s.renderEntity(entity)
	}

	rl.EndMode2D()
	rl.EndTextureMode()
}

func (s *Scene) Key() string {
	return s.key
}

func (s *Scene) RemoveEntity(e common.Entity) {
	s.Components.Viewport.Remove(e)
	s.Components.Color.Remove(e)
	s.Components.Debug.Remove(e)
	s.Components.Fill.Remove(e)
	s.Components.Stroke.Remove(e)
	s.Components.Font.Remove(e)
	s.Components.Interactive.Remove(e)
	s.Components.Opacity.Remove(e)
	s.Components.Origin.Remove(e)
	s.Components.RenderTexture2D.Remove(e)
	s.Components.Size.Remove(e)
	s.Components.SceneGraphNode.Remove(e)
	s.Components.Text.Remove(e)
	s.Components.Texture2D.Remove(e)
	s.Components.Transform.Remove(e)
	s.Components.UUID.Remove(e)

	s.Director.RemoveEntity(e)
}

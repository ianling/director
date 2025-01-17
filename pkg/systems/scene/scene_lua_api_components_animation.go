package scene

import (
	"github.com/gravestench/director/pkg/common"
	"github.com/gravestench/director/pkg/components"
	lua "github.com/yuin/gopher-lua"
)

const (
	luaAnimationComponentName = "animation"
)

/*
example lua:
	anim = scene.components.animation.add(eid)
	anim.frame(2) -- sets frame to 2
	currentFrame = anim.frame()
*/

func (s *Scene) luaExportComponentAnimation(mt *lua.LTable) {
	animationTable := s.Lua.NewTypeMetatable(luaAnimationComponentName)

	s.Lua.SetField(animationTable, "add", s.Lua.NewFunction(s.luaAnimationAdd()))
	s.Lua.SetField(animationTable, "get", s.Lua.NewFunction(s.luaAnimationGet()))
	s.Lua.SetField(animationTable, "remove", s.Lua.NewFunction(s.luaAnimationRemove()))

	s.Lua.SetField(mt, luaAnimationComponentName, animationTable)
}

func (s *Scene) luaAnimationAdd() lua.LGFunction {
	fn := func(L *lua.LState) int {
		if L.GetTop() != 1 {
			return 0
		}

		e := common.Entity(s.Lua.CheckNumber(1))

		animation := s.Components.Animation.Add(e)
		L.Push(s.makeLuaTableComponentAnimation(animation))
		return 1
	}

	return fn
}

func (s *Scene) luaAnimationGet() lua.LGFunction {
	fn := func(L *lua.LState) int {
		if L.GetTop() != 1 {
			return 0
		}

		id := L.CheckNumber(1)
		animation, found := s.Components.Animation.Get(common.Entity(id))

		truthy := lua.LFalse
		if !found {
			L.Push(lua.LNil)
			L.Push(truthy)
			return 2
		} else {
			truthy = lua.LTrue
		}

		table := s.makeLuaTableComponentAnimation(animation)

		L.SetMetatable(table, L.GetTypeMetatable(luaAnimationComponentName))

		L.Push(table)
		L.Push(truthy)

		return 2
	}

	return fn
}

func (s *Scene) luaAnimationRemove() lua.LGFunction {
	fn := func(L *lua.LState) int {
		if L.GetTop() != 1 {
			return 0
		}

		e := common.Entity(s.Lua.CheckNumber(1))

		s.Components.Animation.Remove(e)

		return 0
	}

	return fn
}

func (s *Scene) makeLuaTableComponentAnimation(animation *components.Animation) *lua.LTable {
	table := s.Lua.NewTable()

	s.Lua.SetField(table, "frame", s.Lua.NewFunction(func(L *lua.LState) int {
		if L.GetTop() == 0 {
			L.Push(lua.LNumber(animation.CurrentFrame))

			return 1
		}

		idx := int(L.CheckNumber(1))
		animation.CurrentFrame = idx

		return 0
	}))

	return table
}

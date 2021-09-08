package scene

import (
	"time"

	"github.com/gravestench/akara"
	"github.com/gravestench/director/pkg/common"
)

type imageFactory struct {
	*common.BasicComponents
	common.EntityManager
	cache map[akara.EID]*imageParameters
}

type imageParameters struct {
	uri           string
	width, height int
	debug         bool
}

func (factory *imageFactory) update(s *Scene, _ time.Duration) {
	if !factory.EntityManager.IsInit() {
		factory.EntityManager.Init()
	}

	if factory.cache == nil {
		factory.cache = make(map[akara.EID]*imageParameters)
	}

	factory.EntityManager.ProcessRemovalQueue()
}

func (factory *imageFactory) New(s *Scene, uri string, x, y int) akara.EID {
	e := s.Add.generic.visibleEntity(s)

	trs, _ := s.Components.Transform.Get(e) // this is a component all visible entities have
	trs.Translation.Set(float64(x), float64(y), trs.Translation.Z)

	req := s.Components.FileLoadRequest.Add(e)
	req.Path = uri

	factory.EntityManager.AddEntity(e)

	return e
}

func (factory *imageFactory) putInCache(s *Scene, e akara.EID, uri string, x, y int) {
	entry := &imageParameters{
		uri: uri,
	}

	factory.cache[e] = entry
}

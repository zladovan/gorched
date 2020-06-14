package gorched

import tl "github.com/JoelOtter/termloop"

// SpawnAfter is entity which acts as spawner for another entity.
// It will add Entity to the world right after it will recognize that After entity is not yet in the world.
// After adding Entity to the world it will remove itself from the world.
type SpawnAfter struct {
	Entity, After tl.Drawable
}

// Draw executes spawning logic
func (a *SpawnAfter) Draw(s *tl.Screen) {
	if world, ok := s.Level().(*World); ok {
		for _, e := range world.Entities {
			if e == a.After {
				return
			}
		}
		world.RemoveEntity(a)
		world.AddEntity(a.Entity)
	}

}

// Tick does nothing now
func (a *SpawnAfter) Tick(e tl.Event) {}

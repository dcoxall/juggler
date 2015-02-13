package juggler

import (
	"testing"
)

func TestSpawnerInstanceGeneration(t *testing.T) {
	generated := false
	generator := func(reference string) (Instancer, error) {
		generated = true
		return new(Instance), nil
	}
	spawner := NewSpawner(generator)
	spawner.Fetch("foobar")
	if !generated {
		t.Errorf("Spawner didn't generate new instance with generator function")
	}
}

func TestSpawnerInstanceStorage(t *testing.T) {
	generated := 0
	generator := func(reference string) (Instancer, error) {
		generated += 1
		return new(Instance), nil
	}
	spawner := NewSpawner(generator)
	spawner.Fetch("foobar") // generates a new instance
	spawner.Fetch("barfoo") // generates a new instance
	spawner.Fetch("foobar") // uses previously generated instance
	if generated != 2 {
		t.Errorf("Expected to only generate twice but got %d", generated)
	}
	spawner.Remove("foobar") // forget about foobar
	spawner.Fetch("foobar")  // should regenerate
	if generated != 3 {
		t.Errorf("Expected to regenerate instance")
	}
}

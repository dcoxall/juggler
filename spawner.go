package juggler

import (
	"errors"
	"sync"
)

// Spawner manages a pool of Instancers. The difference between spawner and a
// traditional pool is that the Instancers are not temporary. They are created
// when missing using the generator function but ofuture requests will return
// the same structure.
type Spawner struct {
	generator func(string) (Instancer, error)
	instances map[string]Instancer
	rLocker   sync.Locker
	wLocker   sync.Locker
}

var (
	// Indicates when an Instancer can't be located
	InstanceNotFoundError = errors.New("Unable to locate instance")
)

// NewSpawner will create a new Spawner structure with an appropriate locking
// mechanism and using the provided generator function.
func NewSpawner(gen func(string) (Instancer, error)) *Spawner {
	locker := new(sync.RWMutex)
	return &Spawner{
		generator: gen,
		instances: make(map[string]Instancer),
		rLocker:   locker.RLocker(),
		wLocker:   locker,
	}
}

// Fetch attempts to locate an Instancer that was generated using the provided
// reference otherwise it will pass that reference into the generator function
// returning and storing the result.
func (s *Spawner) Fetch(reference string) (Instancer, error) {
	if i, err := s.getInstance(reference); err != nil {
		if i, err = s.generator(reference); err == nil {
			s.wLocker.Lock()
			defer func(lock sync.Locker) { lock.Unlock() }(s.wLocker)
			s.instances[reference] = i
		}
		return i, err
	} else {
		return i, err
	}
}

// Remove will get rid of any Instancer associated to the provided reference.
func (s *Spawner) Remove(reference string) error {
	if _, err := s.getInstance(reference); err != nil {
		return err
	}
	s.wLocker.Lock()
	defer func(lock sync.Locker) { lock.Unlock() }(s.wLocker)
	delete(s.instances, reference)
	return nil
}

func (s *Spawner) getInstance(reference string) (Instancer, error) {
	s.rLocker.Lock()
	defer func(lock sync.Locker) { lock.Unlock() }(s.rLocker)
	for key, instance := range s.instances {
		if key == reference {
			return instance, nil
		}
	}
	return nil, InstanceNotFoundError
}

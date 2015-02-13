package juggler

import (
	"errors"
	"sync"
)

type Spawner struct {
	generator func(string) (Instancer, error)
	instances map[string]Instancer
	rLocker   sync.Locker
	wLocker   sync.Locker
}

var (
	InstanceNotFoundError = errors.New("Unable to locate instance")
)

func NewSpawner(gen func(string) (Instancer, error)) *Spawner {
	locker := new(sync.RWMutex)
	return &Spawner{
		generator: gen,
		instances: make(map[string]Instancer),
		rLocker:   locker.RLocker(),
		wLocker:   locker,
	}
}

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

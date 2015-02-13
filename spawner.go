package juggler

import "errors"

type Spawner struct {
	generator func(string) (Instancer, error)
	instances map[string]Instancer
}

var (
	InstanceNotFoundError = errors.New("Unable to locate instance")
)

func NewSpawner(gen func(string) (Instancer, error)) *Spawner {
	return &Spawner{
		generator: gen,
		instances: make(map[string]Instancer),
	}
}

func (s *Spawner) Fetch(reference string) (Instancer, error) {
	if i, err := s.getInstance(reference); err != nil {
		if i, err = s.generator(reference); err == nil {
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
	delete(s.instances, reference)
	return nil
}

func (s *Spawner) getInstance(reference string) (Instancer, error) {
	for key, instance := range s.instances {
		if key == reference {
			return instance, nil
		}
	}
	return nil, InstanceNotFoundError
}

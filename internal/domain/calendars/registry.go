package calendars

type Registry struct {
	systems map[string]System
}

func NewRegistry(systems ...System) *Registry {
	registry := &Registry{systems: make(map[string]System, len(systems))}
	for _, system := range systems {
		registry.systems[system.ID()] = system
	}
	return registry
}

func (r *Registry) Get(id string) (System, error) {
	system, ok := r.systems[id]
	if !ok {
		return nil, ErrUnsupportedSystem
	}
	return system, nil
}

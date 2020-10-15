package storage

type Storage struct {
	Engine storable
}

func New(engine storable) *Storage {
	return &Storage{
		Engine: engine,
	}
}

func (s *Storage) PutString(key string, value string) error {
	return s.Engine.put(key, value)
}

func (s *Storage) PutInt(key string, value int) error {
	return s.Engine.put(key, value)
}

func (s Storage) GetString(key string) (string, error) {
	val, err := s.Engine.get(key)
	if err != nil {
		return "", err
	}

	return val.(string), nil
}

func (s Storage) GetInt(key string) (int, error) {
	val, err := s.Engine.get(key)
	if err != nil {
		return 0, err
	}

	return val.(int), nil
}

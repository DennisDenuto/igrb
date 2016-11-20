package store

type Persistor interface {
	Save(key string, val interface{}) error
	Read(key string) (interface{}, error)
}
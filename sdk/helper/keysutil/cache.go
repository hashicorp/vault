package keysutil

type Cache interface {
	Delete(key any)
	Load(key any) (value any, ok bool)
	Store(key, value any)
	Size() int
}

package storeage

type StorageInterface interface {
	IsLocked() (bool, *string)
	Lock(id string) bool
	Unlock(id string) bool
	Delete()
	Put(id string, content []byte) bool
	Get() []byte
}

type StorageBuilderInterface interface {
	Build(project string) (StorageInterface, error)
}

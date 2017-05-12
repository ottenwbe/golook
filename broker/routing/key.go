package routing

import (
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

/*
Key represents the routing key
*/
type Key struct {
	id uuid.UUID
}

/*
NilKey returns a key with a nil value
*/
func NilKey() Key {
	return Key{
		id: uuid.Nil,
	}
}

/*
SysKey returns the key for this golook system
*/
func SysKey() Key {
	u, err := uuid.FromString(golook.GolookSystem.UUID)
	if err != nil {
		log.Error("SysKey() cannot read UUID")
		return NilKey()
	}

	return Key{
		id: u,
	}
}

/*
NewKeyU returns the key with a given uuid
*/
func NewKeyU(key uuid.UUID) Key {
	return Key{
		id: key,
	}
}

/*
NewKey for a given name
*/
func NewKey(name string) Key {
	return Key{
		id: uuid.NewV5(uuid.Nil, name),
	}
}

/*
NewKeyN returns a key for a given namespace and name
*/
func NewKeyN(namespace uuid.UUID, name string) Key {
	return Key{
		id: uuid.NewV5(namespace, name),
	}
}

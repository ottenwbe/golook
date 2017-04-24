package routing

import (
	. "github.com/ottenwbe/golook/broker/runtime"

	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type Key struct {
	id uuid.UUID
}

func NilKey() Key {
	return Key{
		id: uuid.Nil,
	}
}

func SysKey() Key {
	u, err := uuid.FromString(GolookSystem.UUID)
	if err != nil {
		log.Error("SysKey() cannot read UUID")
		return NilKey()
	}

	return Key{
		id: u,
	}
}

func NewKeyU(key uuid.UUID) Key {
	return Key{
		id: key,
	}
}

func NewKey(name string) Key {
	return Key{
		id: uuid.NewV5(uuid.Nil, name),
	}
}

func NewKeyN(namespace uuid.UUID, name string) Key {
	return Key{
		id: uuid.NewV5(namespace, name),
	}
}

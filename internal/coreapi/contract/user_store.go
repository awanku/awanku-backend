package contract

import "github.com/awanku/awanku/pkg/core"

type UserStore interface {
	GetOrCreateByEmail(user *core.User) error
	GetByID(id int64) (*core.User, error)
	Save(user *core.User) error
}

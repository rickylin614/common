package credis

import (
	"context"
	"errors"

	"github.com/go-redsync/redsync/v4"
)

/* init when NewRedisHelper or NewRedisClusterHelper */
var rs *redsync.Redsync

// redsync實現的Lock
func Lock(key string) error {
	if rs == nil {
		return errors.New("there isn't set redsync")
	}
	mut := rs.NewMutex(key)
	return mut.Lock()
}

// redsync實現的LockContext
func LockContext(ctx context.Context, key string) error {
	if rs == nil {
		return errors.New("there isn't set redsync")
	}
	mut := rs.NewMutex(key)
	return mut.LockContext(ctx)
}

// redsync實現的Unlock
func Unlock(key string) (bool, error) {
	if rs == nil {
		return false, errors.New("there isn't set redsync")
	}
	mut := rs.NewMutex(key)
	return mut.Unlock()
}

// redsync實現的UnlockContext
func UnlockContext(ctx context.Context, key string) (bool, error) {
	if rs == nil {
		return false, errors.New("there isn't set redsync")
	}
	mut := rs.NewMutex(key)
	return mut.UnlockContext(ctx)
}

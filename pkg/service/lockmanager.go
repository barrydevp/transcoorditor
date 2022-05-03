package service

import (
	"fmt"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/exception"
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/store"
	"github.com/barrydevp/transcoorditor/pkg/util"
)

var globalLock = util.NewRWLockKey()

func (srv *Service) AcquireLock(key string, owner string, duration time.Duration) (*schema.LockEntry, error) {
	globalLock.Lock(key)
	defer globalLock.Unlock(key)

	existLock, err := srv.s.LockTable().Find(key)
	if err != nil {
		return nil, err
	}
	if existLock != nil && !existLock.IsExpired() {
		return nil, fmt.Errorf("%w: owner(%v)", store.ErrLockExists, existLock.Owner)
	}

	lockEnt := schema.NewLockEntry(key, owner, duration)
	err = srv.s.LockTable().Save(lockEnt)
	if err != nil {
		return nil, exception.Errorf("failed to acquire lock: %w", err)
	}

	return lockEnt, nil
}

func (srv *Service) ReleaseLock(key string, owner string) error {
	lockEnt, err := srv.s.LockTable().FindWithOwner(key, owner)
	if err != nil {
		return exception.Errorf("failed to get lock: %w", err)
	}

	if lockEnt == nil {
		return nil
	}

	return srv.ReleaseLock0(lockEnt)
}

func (srv *Service) ReleaseLock0(lockEnt *schema.LockEntry) error {
	err := srv.s.LockTable().Delete(lockEnt)
	if err != nil {
		return exception.Errorf("failed to release lock: %w", err)
	}

	return nil
}

func (srv *Service) ExtendLock(key string, owner string, duration time.Duration) (*schema.LockEntry, error) {
	globalLock.Lock(key)
	defer globalLock.Unlock(key)

	lockEnt, err := srv.s.LockTable().FindWithOwner(key, owner)

	if err != nil {
		return nil, exception.Errorf("failed to get lock: %w", err)
	}

	if lockEnt.IsExpired() {
		return nil, exception.Errorf("failed to extend lock: %w", store.ErrLockExpired)
	}

	if lockEnt != nil {
		return nil, exception.Errorf("failed to extend lock: %w", store.ErrLockNotFound)
	}

	lockEnt.Extend(duration)
	err = srv.s.LockTable().Update(lockEnt)
	if err != nil {
		return lockEnt, exception.Errorf("failed to extend lock: %w", err)
	}

	return lockEnt, nil
}

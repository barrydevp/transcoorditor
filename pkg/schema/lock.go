package schema

import (
	"time"
	// "github.com/google/uuid"
)

type LockEntry struct {
	Key   string `json:"key" bson:"key"`
	Owner string `json:"owner" bson:"owner"`
	// Uid       string     `json:"uid" bson:"uid"`
	ExpiredAt *time.Time `json:"expiredAt" bson:"expiredAt"`
}

func NewLockEntry(key string, owner string, duration time.Duration) *LockEntry {
	expiredAt := time.Now().Add(duration)

	return &LockEntry{
		Key:       key,
		Owner:     owner,
		ExpiredAt: &expiredAt,
		// Uid:       uuid.NewString(),
	}
}

func (l *LockEntry) IsExpired() bool {
	return l.ExpiredAt.Before(time.Now())
}

func (l *LockEntry) Extend(duration time.Duration) {
	l.ExpiredAt.Add(duration)
}

package controller

import (
	"errors"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/controlplane/reconciler"
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/service"
)

type TimeoutSessionEntry struct {
	TimedoutAt time.Time
	SessionId  string
}

func (en *TimeoutSessionEntry) ExpiredAt() *time.Time {
	return &en.TimedoutAt
}

func shouldRetryAfter2Mins(session *schema.Session, err error) bool {
	if errors.Is(err, service.ErrSessionNotFound) || session == nil {
		// dont retry not found entry
		return false
	}

	if errors.Is(err, service.ErrSessionMaximumRetry) || session.IsMaximumRetry() {
		// dont retry maximum retry entry
		return false
	}

	return true
}

func isSessionNotExpiredYet(session *schema.Session, err error) bool {
	return errors.Is(err, service.ErrSessionNotExpiredYet) && session != nil
}

func (ctrl *Controller) HandleTimeoutSessionRecl(entries []reconciler.ScheduleEntry) []reconciler.ScheduleEntry {
	var newEntries []reconciler.ScheduleEntry
	// use same Now() for all new Entry which reduce the number of loop in schedule
	now := time.Now()

	for _, en := range entries {
		if entry, ok := en.(*TimeoutSessionEntry); ok {
			if session, err := ctrl.srv.TerminateSession(entry.SessionId); err != nil {
				if isSessionNotExpiredYet(session, err) {
					// session timeout has been extended
					newEntries = append(newEntries, &TimeoutSessionEntry{
						TimedoutAt: session.TimedoutAt(),
						SessionId:  entry.SessionId,
					})

				} else if shouldRetryAfter2Mins(session, err) {
					logger.Info("terminate failed, retry after 2 mins")
					// retry after 2 mins

					newEntries = append(newEntries, &TimeoutSessionEntry{
						TimedoutAt: now.Add(2 * time.Minute),
						SessionId:  entry.SessionId,
					})
				}
			}
		} else {
			logger.Error("handleTimeoutSession received malformed entry")
		}
	}

	logger.Info("handleTimeoutSession done!")

	return newEntries
}

func (ctrl *Controller) InitTimeoutSessionQueueRecl() []reconciler.ScheduleEntry {
	var newEntries []reconciler.ScheduleEntry

	sessions, err := ctrl.srv.GetAllUnFinishedSession()
	if err != nil {
		logger.Errorf("Cannot init timeout session queue reconiler")
	}

	for _, session := range sessions {
		newEntries = append(newEntries, &TimeoutSessionEntry{
			TimedoutAt: session.TimedoutAt(),
			SessionId:  session.Id,
		})
	}

	return newEntries
}

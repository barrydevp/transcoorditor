package controller

import (
	"errors"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/controlplane/reconciler"
	"github.com/barrydevp/transcoorditor/pkg/service"
)

type TimeoutSessionEntry struct {
	TimedoutAt time.Time
	SessionId  string
}

func (en *TimeoutSessionEntry) ExpiredAt() *time.Time {
	return &en.TimedoutAt
}

func (ctrl *Controller) HandleTimeoutSessionRecl(entries []reconciler.ScheduleEntry) []reconciler.ScheduleEntry {
	var newEntries []reconciler.ScheduleEntry

	for _, en := range entries {
		if entry, ok := en.(*TimeoutSessionEntry); ok {
			if _, err := ctrl.srv.TerminateSession(entry.SessionId); err != nil {
				if !errors.Is(err, service.ErrSessionNotFound) {
					// retry after 2 mins

					newEntries = append(newEntries, &TimeoutSessionEntry{
						TimedoutAt: time.Now().Add(2 * time.Minute),
						SessionId:  entry.SessionId,
					})
				}
			}
		} else {
			logger.Error("handleTimeoutSession received malformed entry")
		}
	}

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

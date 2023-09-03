package storagescheduler

import (
	syncsched "github.com/M0WA/synced-scheduler/pkgs/syncedscheduler"
)

type schedulerOptionsImpl struct {
	keys      []AssetKey
	isExclude bool
}

// NewSchedulerOptions returns new storage scheduler options
func NewSchedulerOptions(keys []AssetKey, isExclude bool) SchedulerOptions {
	return &schedulerOptionsImpl{keys: keys, isExclude: isExclude}
}

func (o *schedulerOptionsImpl) Keys(m map[AssetKey]Asset) []AssetKey {
	rc := []AssetKey{}

	for ak := range m {
		matched := false
		for _, akk := range o.keys {
			if akk == ak {
				matched = true
				break
			}
		}
		if len(o.keys) == 0 ||
			o.isExclude && !matched ||
			!o.isExclude && matched {
			rc = append(rc, ak)
		}
	}

	return rc
}

// SchedulerFunc is a scheduler function for storage schedulers
type SchedulerFunc func(Resource, SchedulerOptions, map[AssetKey]Asset) (Reservation, error)

type schedulerImpl struct {
	syncsched.Scheduler[AssetKey, Asset, ResourceKey, Resource, Reservation, SchedulerOptions, ResourceReleaser]
	schedFunc SchedulerFunc
	opts      SchedulerOptions
}

// NewScheduler returns a new storage scheduler
func NewScheduler(schedFunc SchedulerFunc) Scheduler {
	return &schedulerImpl{
		Scheduler: syncsched.NewSyncedScheduler[AssetKey, Asset, ResourceKey, Resource, Reservation, SchedulerOptions, ResourceReleaser](),
		schedFunc: schedFunc,
		opts:      NewSchedulerOptions([]AssetKey{}, false),
	}
}

func (s *schedulerImpl) Schedule(r Resource) (Reservation, error) {
	return s.Scheduler.ScheduleResourceLocked(r, s.opts, s.schedFunc)
}

func (s *schedulerImpl) Remove(r Reservation) error {
	err := s.Scheduler.RemoveResource(r, newResourceReleaser())
	return err
}

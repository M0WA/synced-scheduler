package syncedscheduler

import (
	"sync"
)

type syncedScheduler[
	assetKey AssetKey,
	asset Asset[assetKey, resourceKey, resource],
	resourceKey ResourceKey,
	resource Resource[resourceKey],
	reservation Reservation[assetKey, resourceKey, resource],
	schedOpts SchedulerOptions,
	resourceReleaser ResourceReleaser[assetKey, asset, resourceKey, resource]] struct {
	m map[assetKey]asset
	l sync.Mutex
}

func (c *syncedScheduler[assetKey, asset, resourceKey, resource, reservation, schedOpts, resourceReleaser]) AddAsset(v asset) error {
	ak := v.AssetKey()

	c.l.Lock()
	defer c.l.Unlock()

	if _, ok := c.m[ak]; ok {
		return ErrorAssetAlreadyExists
	}
	c.m[ak] = v
	return nil
}

func (c *syncedScheduler[assetKey, asset, resourceKey, resource, reservation, schedOpts, resourceReleaser]) RemoveAsset(k assetKey) error {
	c.l.Lock()
	defer c.l.Unlock()

	if _, ok := c.m[k]; !ok {
		return ErrorAssetNotExists
	}

	delete(c.m, k)
	return nil
}

func (c *syncedScheduler[assetKey, asset, resourceKey, resource, reservation, schedOpts, resourceReleaser]) ScheduleResourceLocked(r resource, o schedOpts, schedFunc func(resource, schedOpts, map[assetKey]asset) (reservation, error)) (reservation, error) {
	c.l.Lock()
	defer c.l.Unlock()

	res, err := schedFunc(r, o, c.m)
	return res, err
}

func (c *syncedScheduler[assetKey, asset, resourceKey, resource, reservation, schedOpts, resourceReleaser]) ScheduleResource(r resource, o schedOpts, schedFunc func(resource, schedOpts, *sync.Mutex, map[assetKey]asset) (reservation, error)) (reservation, error) {
	res, err := schedFunc(r, o, &c.l, c.m)
	return res, err
}

func (c *syncedScheduler[assetKey, asset, resourceKey, resource, reservation, schedOpts, resourceReleaser]) RemoveResource(resv reservation, releaser resourceReleaser) error {
	c.l.Lock()
	defer c.l.Unlock()

	ass, ok := c.m[resv.AssetKey()]
	if !ok {
		return ErrorAssetNotExists
	}
	res := ass.Resources()
	_, ok = res[resv.Resource().ResourceKey()]
	if !ok {
		return ErrorResourceNotExists
	}
	if err := releaser.ReleaseResource(ass, resv.Resource()); err != nil {
		return err
	}
	delete(res, resv.Resource().ResourceKey())
	return nil
}

// NewSyncedScheduler creates a synced scheduler
func NewSyncedScheduler[
	assetKey AssetKey,
	asset Asset[assetKey, resourceKey, resource],
	resourceKey ResourceKey,
	resource Resource[resourceKey],
	reservation Reservation[assetKey, resourceKey, resource],
	schedOpts SchedulerOptions,
	resourceReleaser ResourceReleaser[assetKey, asset, resourceKey, resource]]() Scheduler[assetKey, asset, resourceKey, resource, reservation, schedOpts, resourceReleaser] {
	return &syncedScheduler[assetKey, asset, resourceKey, resource, reservation, schedOpts, resourceReleaser]{m: make(map[assetKey]asset)}
}

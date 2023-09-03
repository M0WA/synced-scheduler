package storagescheduler

import (
	syncsched "github.com/M0WA/synced-scheduler/pkgs/syncedscheduler"
)

type assetImpl struct {
	key AssetKey
	m   map[ResourceKey]Resource

	capacity uint64
	usage    uint64
}

// NewAsset creates a new storage asset
func NewAsset(key AssetKey, capacity uint64) Asset {
	return &assetImpl{key: key, capacity: capacity, m: make(map[ResourceKey]Resource)}
}

func (a *assetImpl) AssetKey() AssetKey {
	return a.key
}

func (a *assetImpl) Resources() map[ResourceKey]Resource {
	return a.m
}

func (a *assetImpl) IncrementUsage(delta uint64) error {
	tmp := a.usage + delta
	if tmp > a.capacity {
		return syncsched.ErrorOutOfCapacity
	}
	a.usage = tmp
	return nil
}

func (a *assetImpl) DecrementUsage(delta uint64) error {
	if delta > a.usage {
		return syncsched.ErrorOutOfCapacity
	}
	a.usage -= delta
	return nil
}

func (a *assetImpl) FillRatio() float64 {
	if a.capacity == 0 {
		return 100.0
	} else {
		return float64(a.usage) / float64(a.capacity)
	}
}

func (a *assetImpl) HasCapacityFree(additional uint64) bool {
	return a.usage+additional < a.capacity
}

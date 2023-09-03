package storagescheduler

import (
	syncsched "github.com/M0WA/synced-scheduler/pkgs/syncedscheduler"
)

type ResourceKey string
type AssetKey string

type Asset interface {
	syncsched.Asset[AssetKey, ResourceKey, Resource]
	IncrementUsage(uint64) error
	DecrementUsage(uint64) error
	FillRatio() float64
	HasCapacityFree(uint64) bool
}

type Resource interface {
	syncsched.Resource[ResourceKey]
	Size() uint64
}

type Reservation syncsched.Reservation[AssetKey, ResourceKey, Resource]

type SchedulerOptions interface {
	Keys(map[AssetKey]Asset) []AssetKey
}

type ResourceReleaser interface {
	syncsched.ResourceReleaser[AssetKey, Asset, ResourceKey, Resource]
}

type Scheduler interface {
	syncsched.Scheduler[AssetKey, Asset, ResourceKey, Resource, Reservation, SchedulerOptions, ResourceReleaser]
	Schedule(Resource) (Reservation, error)
	Remove(Reservation) error
}

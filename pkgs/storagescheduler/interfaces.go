package storagescheduler

import (
	syncsched "github.com/M0WA/synced-scheduler/pkgs/syncedscheduler"
)

// ResourceKey for storage resources
type ResourceKey string

// AssetKey for storage assets
type AssetKey string

// Asset is a storage asset
type Asset interface {
	syncsched.Asset[AssetKey, ResourceKey, Resource]

	// IncrementUsage increments the usage of an asset
	IncrementUsage(uint64) error

	// DecrementUsage decrements the usage of an asset
	DecrementUsage(uint64) error

	// FillRatio of an asset in percent
	FillRatio() float64

	// HasCapacityFree returns true if asset has enough capacity
	HasCapacityFree(uint64) bool
}

// Resource is a storage resource
type Resource interface {
	syncsched.Resource[ResourceKey]

	// Size returns the size of this storage resource
	Size() uint64
}

// Reserveration is a storage reservation
type Reservation syncsched.Reservation[AssetKey, ResourceKey, Resource]

// SchedulerOptions are storage scheduler options
type SchedulerOptions interface {
	// Keys returns the allowed asset keys
	Keys(map[AssetKey]Asset) []AssetKey
}

// ResourceReleaser releases storage resources from assets
type ResourceReleaser interface {
	syncsched.ResourceReleaser[AssetKey, Asset, ResourceKey, Resource]
}

// Scheduler schedules storage resources on storage assets
type Scheduler interface {
	// Schedule schedules a resource and returns a reservation
	Schedule(Resource) (Reservation, error)

	// Remove removes a reservation
	Remove(Reservation) error

	// AddAsset adds an asset for scheduling resources
	AddAsset(Asset) error

	// RemoveAsset removes an asset for scheduling resources
	RemoveAsset(AssetKey) error
}

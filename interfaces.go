package syncedscheduler

import "sync"

// AssetKey is the type placeholder for asset keys
type AssetKey comparable

// Asset is the type placeholder for assets
type Asset[assetKey AssetKey, resourceKey ResourceKey, resource Resource[resourceKey]] interface {
	any
	AssetKey() assetKey
	Resources() map[resourceKey]resource
}

// ResourceKey is the type placeholder for resource keys
type ResourceKey comparable

// Resource is the interface type placeholder for resources
type Resource[resourceKey ResourceKey] interface {
	any
	ResourceKey() resourceKey
}

// Reservation is the interface type placeholder for reservations
type Reservation[assetKey AssetKey, resourceKey ResourceKey, resource Resource[resourceKey]] interface {
	AssetKey() assetKey
	Resource() resource
}

// Scheduler is the interface type placeholder for scheduler implementations
type Scheduler[assetKey AssetKey, asset Asset[assetKey, resourceKey, resource], resourceKey ResourceKey, resource Resource[resourceKey], reservation Reservation[assetKey, resourceKey, resource]] interface {
	// AddAsset adds an asset for scheduling resources
	AddAsset(assetKey, asset) error

	// RemoveAsset removes an asset for scheduling resources
	RemoveAsset(assetKey) error

	// ScheduleResourceLocked will schedule a resource and will lock cache automatically
	ScheduleResourceLocked(resource, func(resource, map[assetKey]asset) (reservation, error)) (reservation, error)

	// ScheduleResource will schedule a resource, cache locking has to be done in the given scheduling function yourself
	ScheduleResource(resource, func(resource, *sync.Mutex, map[assetKey]asset) (reservation, error)) (reservation, error)

	// RemoveResource removes a resource allocation from an asset
	RemoveResource(reservation) error
}

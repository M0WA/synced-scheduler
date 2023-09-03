package syncedscheduler

type reservationImpl[assetKey AssetKey, resourceKey ResourceKey, resource Resource[resourceKey]] struct {
	assetKey assetKey
	resource resource
}

func (r *reservationImpl[assetKey, resourceKey, resource]) AssetKey() assetKey {
	return r.assetKey
}

func (r *reservationImpl[assetKey, resourceKey, resource]) Resource() resource {
	return r.resource
}

// NewReservation will create a Reservation from an asset key and a resource
func NewReservation[assetKey AssetKey, resourceKey ResourceKey, resource Resource[resourceKey]](ak assetKey, res resource) Reservation[assetKey, resourceKey, resource] {
	return &reservationImpl[assetKey, resourceKey, resource]{assetKey: ak, resource: res}
}

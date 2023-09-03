package storagescheduler

type resourceImpl struct {
	key  ResourceKey
	size uint64
}

func NewResource(key ResourceKey, size uint64) Resource {
	return &resourceImpl{key: key, size: size}
}

func (r *resourceImpl) ResourceKey() ResourceKey {
	return r.key
}

func (r *resourceImpl) Size() uint64 {
	return r.size
}

type resourceReleaserImpl struct{}

func newResourceReleaser() ResourceReleaser {
	return &resourceReleaserImpl{}
}

func (*resourceReleaserImpl) ReleaseResource(a Asset, r Resource) error {
	a.DecrementUsage(r.Size())
	return nil
}

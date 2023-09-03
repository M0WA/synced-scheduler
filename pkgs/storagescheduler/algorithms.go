package storagescheduler

import (
	syncsched "github.com/M0WA/synced-scheduler/pkgs/syncedscheduler"
)

var LowestFillRatioAlgorithm SchedulerFunc = func(r Resource, o SchedulerOptions, m map[AssetKey]Asset) (Reservation, error) {
	aks := o.Keys(m)
	var minFill float64 = 100.0
	var akp *AssetKey = nil
	for _, ak := range aks {
		ass := m[ak]
		if !ass.HasCapacityFree(r.Size()) {
			continue
		}
		fr := ass.FillRatio()
		if fr < minFill {
			minFill = fr
			akp = &ak
		}
	}
	if akp == nil {
		return nil, syncsched.ErrorOutOfCapacity
	}
	ak := *akp
	as := m[ak]
	if err := as.IncrementUsage(r.Size()); err != nil {
		return nil, err
	}
	as.Resources()[r.ResourceKey()] = r
	return syncsched.NewReservation[AssetKey, ResourceKey, Resource](*akp, r), nil
}

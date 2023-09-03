package syncedscheduler_test

import (
	"crypto/rand"
	"fmt"
	"sync"
	"testing"

	sched "github.com/M0WA/synced-scheduler"
)

func makeUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func newTestScheduler() sched.Scheduler[testAssetKey, testAsset, testResourceKey, testResource, sched.Reservation[testAssetKey, testResourceKey, testResource]] {
	return sched.NewSyncedScheduler[testAssetKey, testAsset, testResourceKey, testResource, sched.Reservation[testAssetKey, testResourceKey, testResource]]()
}

func TestCache(t *testing.T) {
	c := newTestScheduler()

	a1 := newTestAsset(testAssetKey(makeUUID()), 1)
	res1 := newTestResource(testResourceKey(makeUUID()))

	t.Run("add asset", func(t *testing.T) {
		if err := c.AddAsset(a1.AssetKey(), a1); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("prevent re-add asset", func(t *testing.T) {
		if err := c.AddAsset(a1.AssetKey(), a1); err != sched.ErrorResourceAlreadyExists {
			t.Fatal(err)
		}
	})

	t.Run("schedule resource locked", func(t *testing.T) {
		if _, err := c.ScheduleResourceLocked(res1, func(tr testResource, m map[testAssetKey]testAsset) (sched.Reservation[testAssetKey, testResourceKey, testResource], error) {
			for k := range m {
				m[k].Resources()[tr.ResourceKey()] = tr
				return sched.NewReservation[testAssetKey, testResourceKey, testResource](m[k].AssetKey(), tr), nil
			}
			return nil, sched.ErrorOutOfCapacity
		}); err != nil {
			t.Fatal(err)
		}
	})

	var resv sched.Reservation[testAssetKey, testResourceKey, testResource]
	t.Run("schedule resource unlocked", func(t *testing.T) {
		var err error
		if resv, err = c.ScheduleResource(res1, func(tr testResource, l *sync.Mutex, m map[testAssetKey]testAsset) (sched.Reservation[testAssetKey, testResourceKey, testResource], error) {
			l.Lock()
			defer l.Unlock()

			for k := range m {
				m[k].Resources()[tr.ResourceKey()] = tr
				return sched.NewReservation[testAssetKey, testResourceKey, testResource](m[k].AssetKey(), tr), nil
			}
			return nil, sched.ErrorOutOfCapacity
		}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("remove resource", func(t *testing.T) {
		if err := c.RemoveResource(resv); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("remove asset", func(t *testing.T) {
		if err := c.RemoveAsset(a1.AssetKey()); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("prevent remove of invalid asset", func(t *testing.T) {
		if err := c.RemoveAsset(a1.AssetKey()); err != sched.ErrorResourceNotExists {
			t.Fatal(err)
		}
	})
}

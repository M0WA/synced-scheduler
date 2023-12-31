package syncedscheduler_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	sched "github.com/M0WA/synced-scheduler/pkgs/syncedscheduler"
)

func testBenchInitAssets(s testScheduler, count int) ([]testAssetKey, error) {
	var wg sync.WaitGroup
	wg.Add(count)
	c := make(chan testAssetKey)

	go func() {
		wg.Wait()
		close(c)
	}()

	rc := []testAssetKey{}
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			ak := testAssetKey(makeUUID())
			if err := s.AddAsset(newTestAsset(ak, 100)); err != nil {
				panic(err)
			}
			c <- ak
		}()
	}

	for ak := range c {
		rc = append(rc, ak)
	}
	return rc, nil
}

func testBenchScheduleFunc(tr testResource, o testSchedulerOpts, m map[testAssetKey]testAsset) (testReservation, error) {
	keys := []testAssetKey{}
	for k := range m {
		keys = append(keys, k)
	}
	if len(keys) < 1 {
		return nil, sched.ErrorOutOfCapacity
	}
	pos := 0
	if len(keys) < 1 {
		pos = rand.Intn(len(keys) - 1)
	}
	kk := keys[pos]
	mak := m[kk]
	ak := mak.AssetKey()
	resm := mak.Resources()
	resm[tr.ResourceKey()] = tr
	return sched.NewReservation[testAssetKey, testResourceKey, testResource](ak, tr), nil
}

func testBenchScheduleResources(s testScheduler, count int, benchFunc benchFuncType) ([]testReservation, error) {
	var wg sync.WaitGroup
	wg.Add(count)
	c := make(chan testReservation)

	go func() {
		wg.Wait()
		close(c)
	}()

	resrvs := []testReservation{}
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			rr, err := s.ScheduleResourceLocked(newTestResource(testResourceKey(makeUUID())), nil, benchFunc)
			if err != nil {
				panic(err)
			}
			c <- rr
		}()
	}

	for rr := range c {
		resrvs = append(resrvs, rr)
	}
	return resrvs, nil
}

func testBenchRemoveResources(s testScheduler, resrvs []testReservation) error {
	var wg sync.WaitGroup
	wg.Add(len(resrvs))
	for _, r := range resrvs {
		go func(rrr testReservation) {
			defer wg.Done()
			if err := s.RemoveResource(rrr, newTestResourceReleaser()); err != nil {
				panic(err)
			}
		}(r)
	}
	wg.Wait()
	return nil
}

func testBenchRemoveAssets(s testScheduler, aks []testAssetKey) error {
	var wg sync.WaitGroup
	wg.Add(len(aks))
	for _, a := range aks {
		go func(aaa testAssetKey) {
			defer wg.Done()
			if err := s.RemoveAsset(aaa); err != nil {
				panic(err)
			}
		}(a)
	}
	wg.Wait()
	return nil
}

type benchFuncType func(testResource, testSchedulerOpts, map[testAssetKey]testAsset) (testReservation, error)

func benchSchedulerRun(s testScheduler, assetCount int, resourceCount int, benchFunc benchFuncType) error {
	aks, err := testBenchInitAssets(s, assetCount)
	if err != nil {
		return err
	}

	resrvs, err := testBenchScheduleResources(s, resourceCount, benchFunc)
	if err != nil {
		return err
	}
	if err := testBenchRemoveResources(s, resrvs); err != nil {
		return err
	}
	if err := testBenchRemoveAssets(s, aks); err != nil {
		return err
	}

	return nil
}

func BenchmarkSyncedScheduler(b *testing.B) {
	assetSteps := []int{10, 100, 1500}
	resourceSteps := []int{10, 1000, 15000, 150000}
	benchfuncs := map[string]benchFuncType{"testBenchScheduleFunc": testBenchScheduleFunc}

	for benchFuncName, benchFunc := range benchfuncs {
		b.Run(benchFuncName, func(b *testing.B) {
			for _, assetCount := range assetSteps {
				b.Run(fmt.Sprintf("%d assets", assetCount), func(b *testing.B) {
					for _, resourceCount := range resourceSteps {
						b.Run(fmt.Sprintf("%d resources", resourceCount), func(b *testing.B) {
							if err := benchSchedulerRun(newTestSyncedScheduler(), assetCount, resourceCount, benchFunc); err != nil {
								b.Fatal(err)
							}
						})
					}
				})
			}
		})
	}
}

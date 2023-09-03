package storagescheduler_test

import (
	"crypto/rand"
	"fmt"
	"sync"
	"testing"

	stgsched "github.com/M0WA/synced-scheduler/pkgs/storagescheduler"
)

const (
	MegaByte = 1024 * 1024
)

const (
	assetCapacity uint64 = 10000000000000 * MegaByte
	resourceUsage uint64 = 1 * MegaByte
)

func makeUUID() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func testBenchInitAssets(s stgsched.Scheduler, count int, capacity uint64) ([]stgsched.AssetKey, error) {
	var wg sync.WaitGroup
	wg.Add(count)
	c := make(chan stgsched.AssetKey)

	go func() {
		wg.Wait()
		close(c)
	}()

	rc := []stgsched.AssetKey{}
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			ak := stgsched.AssetKey(makeUUID())
			if err := s.AddAsset(stgsched.NewAsset(ak, capacity)); err != nil {
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

func testBenchScheduleResources(s stgsched.Scheduler, count int) ([]stgsched.Reservation, error) {
	var wg sync.WaitGroup
	wg.Add(count)
	c := make(chan stgsched.Reservation)

	go func() {
		wg.Wait()
		close(c)
	}()

	resrvs := []stgsched.Reservation{}
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			tr := stgsched.NewResource(stgsched.ResourceKey(makeUUID()), resourceUsage)
			rr, err := s.Schedule(tr)
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

func testBenchRemoveResources(s stgsched.Scheduler, resrvs []stgsched.Reservation) error {
	var wg sync.WaitGroup
	wg.Add(len(resrvs))
	for _, r := range resrvs {
		var rrr stgsched.Reservation = r
		go func() {
			defer wg.Done()
			if err := s.Remove(rrr); err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
	return nil
}

func testBenchRemoveAssets(s stgsched.Scheduler, aks []stgsched.AssetKey) error {
	var wg sync.WaitGroup
	wg.Add(len(aks))
	for _, a := range aks {
		var aaa stgsched.AssetKey = a
		go func() {
			defer wg.Done()
			if err := s.RemoveAsset(aaa); err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
	return nil
}

func benchSchedulerRun(s stgsched.Scheduler, assetCount int, resourceCount int) error {

	aks, err := testBenchInitAssets(s, assetCount, assetCapacity)
	if err != nil {
		return err
	}

	resrvs, err := testBenchScheduleResources(s, resourceCount)
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

func BenchmarkStorageScheduler(b *testing.B) {
	assetSteps := []int{10, 100, 1500}
	resourceSteps := []int{10, 1000, 15000, 150000}
	benchfuncs := map[string]stgsched.SchedulerFunc{"LowestFillRatioAlgorithm": stgsched.LowestFillRatioAlgorithm}

	for benchFuncName, benchFunc := range benchfuncs {
		b.Run(benchFuncName, func(b *testing.B) {
			for _, assetCount := range assetSteps {
				b.Run(fmt.Sprintf("%d assets", assetCount), func(b *testing.B) {
					for _, resourceCount := range resourceSteps {
						b.Run(fmt.Sprintf("%d resources", resourceCount), func(b *testing.B) {
							s := stgsched.NewScheduler(benchFunc)
							if err := benchSchedulerRun(s, assetCount, resourceCount); err != nil {
								b.Fatal(err)
							}
						})
					}
				})
			}
		})
	}
}

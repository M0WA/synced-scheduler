package syncedscheduler_test

import (
	sched "github.com/M0WA/synced-scheduler"
)

// user defined key for cache
type testAssetKey string

// user defined interface for cache contents
type testAsset interface {
	sched.Asset[testAssetKey, testResourceKey, testResource]

	// some example properties
	Usage() int
	SetUsage(int)
}

type testAssetImpl struct {
	usage     int
	id        testAssetKey
	resources map[testResourceKey]testResource
}

func (a *testAssetImpl) AssetKey() testAssetKey {
	return a.id
}

func (a *testAssetImpl) Resources() map[testResourceKey]testResource {
	return a.resources
}

func (a *testAssetImpl) Usage() int {
	return a.usage
}

func (a *testAssetImpl) SetUsage(usage int) {
	a.usage = usage
}

func newTestAsset(id testAssetKey, usage int) testAsset {
	return &testAssetImpl{usage: usage, id: id, resources: map[testResourceKey]testResource{}}
}

package syncedscheduler_test

import (
	sched "github.com/M0WA/synced-scheduler/pkgs/syncedscheduler"
)

type testResourceKey string

type testResource interface {
	sched.Resource[testResourceKey]
	GetUsage() int
	SetUsage(int) error
}

type testResourceImpl struct {
	usage int
	key   testResourceKey
}

func (t *testResourceImpl) GetUsage() int {
	return t.usage
}

func (t *testResourceImpl) SetUsage(usage int) error {
	t.usage = usage
	return nil
}

func (t *testResourceImpl) ResourceKey() testResourceKey {
	return t.key
}

func newTestResource(id testResourceKey) testResource {
	return &testResourceImpl{key: id}
}

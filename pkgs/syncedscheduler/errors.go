package syncedscheduler

import "errors"

// ErrorResourceNotExists resource does not exist
var ErrorResourceNotExists error = errors.New("resource does not exist")

// ErrorAssetAlreadyExists assetalready exists
var ErrorAssetAlreadyExists error = errors.New("asset already exists")

// ErrorAssetNotExists asset does not exist
var ErrorAssetNotExists error = errors.New("asset does not exist")

// ErrorOutOfCapacity out of capacity, cannot schedule
var ErrorOutOfCapacity error = errors.New("out of capacity")

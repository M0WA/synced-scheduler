package syncedscheduler

import "errors"

// ErrorResourceAlreadyExists resource already exists
var ErrorResourceAlreadyExists error = errors.New("already exists")

// ErrorResourceNotExists resource does not exist
var ErrorResourceNotExists error = errors.New("not exists")

// ErrorOutOfCapacity out of capacity, cannot schedule
var ErrorOutOfCapacity error = errors.New("out of capacity")

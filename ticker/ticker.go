package ticker

import (
	"time"
)

func TickTasks(taskFn func(), duration time.Duration) {
	tick := time.NewTimer(duration)
	defer tick.Stop()

	start := time.Now()
	for {
		taskFn()

		cost := time.Since(start)
		ttl := duration - cost
		if ttl < 0 {
			ttl = 0
		}
		tick.Reset(ttl)
		select {
		case <-tick.C:
		}
		start = time.Now()
	}
}

func TickTasksWithValues(taskFn func(values ...interface{}), duration time.Duration, values ...interface{}) {
	tick := time.NewTimer(duration)
	defer tick.Stop()

	start := time.Now()
	for {
		taskFn(values...)

		cost := time.Since(start)
		ttl := duration - cost
		if ttl < 0 {
			ttl = 0
		}
		tick.Reset(ttl)
		select {
		case <-tick.C:
		}
		start = time.Now()
	}
}

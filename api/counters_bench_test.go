package api

// go test -v ./... -bench=. -run=^$ -benchmem --count 5
import (
	"testing"
)

var counters = NewAllCounters()

type counterOp struct {
	name      string
	increment func()
	get       func() uint64
}

var counterOps = []counterOp{
	{
		name:      "BadAF",
		increment: func() { counters.badAF.Increment() },
		get:       func() uint64 { return counters.badAF.Get() },
	},
	{
		name:      "ScopedValue",
		increment: func() { counters.scopedValue() },
		get:       func() uint64 { return counters.scopedValue() },
	},
	{
		name:      "Mutex",
		increment: func() { counters.mutex.Increment() },
		get:       func() uint64 { return counters.mutex.Get() },
	},
	{
		name:      "Atomic",
		increment: func() { counters.atmoic.Increment() },
		get:       func() uint64 { return counters.atmoic.Get() },
	},
	{
		name:      "Semaphore",
		increment: func() { counters.semaphore.Increment() },
		get:       func() uint64 { return counters.semaphore.Get() },
	},
	{
		name:      "Channel",
		increment: func() { counters.channel.IncrementAndGet() },
		get:       func() uint64 { return counters.channel.IncrementAndGet() },
	},
}

func BenchmarkCounters(b *testing.B) {
	b.Run("Increment", func(b *testing.B) {
		for _, op := range counterOps {
			b.Run(op.name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					op.increment()
				}
			})
		}
	})

	b.Run("Get", func(b *testing.B) {
		for _, op := range counterOps {
			b.Run(op.name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = op.get()
				}
			})
		}
	})
}

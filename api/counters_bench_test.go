package api

// go test -v ./... -bench=. -run=^$ -benchmem --count 5
import (
	"testing"
)

// var counters = NewAllCounters()

// // func BenchmarkChannelWriteRead(b *testing.B) {
// // 	// go counters.channel.run()
// // 	for i := 0; i < b.N; i++ {
// // 		counters.channel.IncrementAndGet()
// // 	}
// // }

// func BenchmarkSemaphoreWrite(b *testing.B) {

// 	ctx := context.Background()
// 	for i := 0; i < b.N; i++ {
// 		if err := counters.semaphore.Acquire(ctx, 1); err != nil {
// 			b.Fatal(err)
// 		}
// 		counters.semaphore.value++
// 		counters.semaphore.Release(1)
// 	}
// }

// func BenchmarkSemaphoreRead(b *testing.B) {
// 	ctx := context.Background()
// 	for i := 0; i < b.N; i++ {
// 		if err := counters.semaphore.Acquire(ctx, 1); err != nil {
// 			b.Fatal(err)
// 		}
// 		_ = counters.semaphore.value
// 		counters.semaphore.Release(1)
// 	}
// }

// func BenchmarkAtomicWrite(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		counters.atmoic.value.Add(1)
// 	}
// }

// func BenchmarkMutexWrite(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		counters.mutex.Lock()
// 		counters.mutex.value++
// 		counters.mutex.Unlock()
// 	}
// }

// func BenchmarkAtomicRead(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		_ = counters.atmoic.value.Load()
// 	}
// }

// func BenchmarkMutexRead(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		counters.mutex.Lock()
// 		_ = counters.mutex.value
// 		counters.mutex.Unlock()
// 	}
// }

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
		increment: func() { counters.atmoic.value.Add(1) },
		get:       func() uint64 { return counters.atmoic.value.Load() },
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

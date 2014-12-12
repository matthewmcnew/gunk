package faketimeprovider_test

import (
	"fmt"
	"time"

	"github.com/matthewmcnew/gunk/timeprovider/faketimeprovider"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FakeTimeProvider", func() {
	const Δ time.Duration = 10 * time.Millisecond

	var (
		timeProvider *faketimeprovider.FakeTimeProvider
		initialTime  time.Time
	)

	BeforeEach(func() {
		initialTime = time.Date(2014, 1, 1, 3, 0, 30, 0, time.UTC)
		timeProvider = faketimeprovider.New(initialTime)
	})

	Describe("Now", func() {
		It("returns the current time, w/o race conditions", func() {
			go timeProvider.Increment(time.Minute)
			Eventually(timeProvider.Now).Should(Equal(initialTime.Add(time.Minute)))
		})
	})

	Describe("Sleep", func() {
		It("blocks until the given interval elapses", func() {
			doneSleeping := make(chan struct{})
			go func() {
				timeProvider.Sleep(10 * time.Second)
				close(doneSleeping)
				fmt.Println("done sleeping")
			}()

			Consistently(doneSleeping, Δ).ShouldNot(BeClosed())

			timeProvider.Increment(5 * time.Second)
			Consistently(doneSleeping, Δ).ShouldNot(BeClosed())

			timeProvider.Increment(4 * time.Second)
			Consistently(doneSleeping, Δ).ShouldNot(BeClosed())

			timeProvider.Increment(1 * time.Second)
			Eventually(doneSleeping).Should(BeClosed())
		})

		FIt("should return in order", func() {
			order := &[]int{}
			go func() {
				fmt.Println("begin the sleep")
				timeProvider.Sleep(10 * time.Second)
				*order = append(*order, 10)
				fmt.Println(order)
				fmt.Println("Firing")
			}()

			go func() {
				timeProvider.Sleep(10 * time.Second)
				*order = append(*order, 10)
			}()

			go func() {
				timeProvider.Sleep(15 * time.Second)
				*order = append(*order, 15)
			}()

			go func() {
				fmt.Println("begin the sleep")
				timeProvider.Sleep(5 * time.Second)
				*order = append(*order, 5)
				fmt.Println(order)
			}()

			Consistently(*order, 1).Should(Equal([]int{}))

			timeProvider.Increment(25 * time.Second)
				time.Sleep(10*time.Second)
				fmt.Println(*order)
			Eventually(*order).Should(Equal([]int{5, 10, 15, 20}))
		})
	})

	Describe("WatcherCount", func() {
		Context("when a timer is created", func() {
			It("increments the watcher count", func() {
				timeProvider.NewTimer(time.Second)
				Ω(timeProvider.WatcherCount()).Should(Equal(1))

				timeProvider.NewTimer(2 * time.Second)
				Ω(timeProvider.WatcherCount()).Should(Equal(2))
			})
		})

		//		Context("when a timer fires", func() {
		//			It("increments the watcher count", func() {
		//				timeProvider.NewTimer(time.Second)
		//				Ω(timeProvider.WatcherCount()).Should(Equal(1))
		//
		//				timeProvider.Increment(time.Second)
		//				Ω(timeProvider.WatcherCount()).Should(Equal(0))
		//			})
		//		})
	})
})

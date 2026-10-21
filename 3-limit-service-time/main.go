//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"context"
	"sync/atomic"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	ctx, cancel := context.Background(), context.CancelFunc(nil)
	if u.IsPremium == false {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(10-atomic.LoadInt64(&u.TimeUsed))*time.Second)
		defer cancel()
	}
	done := make(chan struct{})
	processDuration := int64(0)
	ticker := time.NewTicker(time.Second)
	go func() {
		process()
		done <- struct{}{}
	}()
	ticker.Stop()
	for {
		select {
		case <-ticker.C:
			processDuration++
		case <-done:
			atomic.StoreInt64(&u.TimeUsed, atomic.LoadInt64(&u.TimeUsed)+processDuration)
			return true
		case <-ctx.Done():
			atomic.StoreInt64(&u.TimeUsed, 10)
			return false
		}
	}

	/*
		if u.IsPremium == false {
			timeout := time.Second * time.Duration(10-atomic.LoadInt64(&u.TimeUsed))
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			doneInSeconds := make(chan int64)
			go func(ctx context.Context, cancel context.CancelFunc) {
				start := time.Now()
				process()
				doneInSeconds <- int64(time.Since(start).Seconds())
			}(ctx, cancel)
			for {
				select {
				case processDuration := <-doneInSeconds:
					atomic.StoreInt64(&u.TimeUsed, atomic.LoadInt64(&u.TimeUsed)+processDuration)
					return true
				case <-ctx.Done():
					atomic.StoreInt64(&u.TimeUsed, 10)
					return false
				}
			}
		} else {
			doneInSeconds := make(chan int64)
			go func() {
				start := time.Now()
				process()
				doneInSeconds <- int64(time.Since(start).Seconds())
			}()
			atomic.StoreInt64(&u.TimeUsed, atomic.LoadInt64(&u.TimeUsed)+<-doneInSeconds)
			return true
		}
	*/
}

func main() {
	RunMockServer()
}

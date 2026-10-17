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
	if u.TimeUsed >= 10 {
		return false
	}
	ticker := time.NewTicker(time.Second)
	done := make(chan struct{})
	go func() {
		process()
		done <- struct{}{}
	}()
	for {
		select {
		case <-ticker.C:
			atomic.StoreInt64(&u.TimeUsed, atomic.AddInt64(&u.TimeUsed, 1))
			if !u.IsPremium && atomic.LoadInt64(&u.TimeUsed) >= 10 {
				return false
			}
		case <-done:
			return true
		}
	}
}

func main() {
	RunMockServer()
}

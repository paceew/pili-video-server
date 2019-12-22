package main

import "log"

type ConnLimiter struct {
	concurrentConn int
	bucket         chan int
}

func NewConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		concurrentConn: cc,
		bucket:         make(chan int, cc),
	}
}

func (cl *ConnLimiter) GetConn() bool {
	if len(cl.bucket) >= cl.concurrentConn {
		log.Printf("Reached the rate limitation.\n")
		return false
	}

	cl.bucket <- 1
	log.Printf("get a connect!\n")
	return true
}

func (cl *ConnLimiter) ReleaseConn() {
	_ = <-cl.bucket
	log.Printf("release a connect!\n")
}

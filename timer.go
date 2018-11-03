package main

import (
	"time"
)

func Tick() {
	c := time.Tick(1 * time.Second)
	for {
		<-c
		UpdateState()
	}
}

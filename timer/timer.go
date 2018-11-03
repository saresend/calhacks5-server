package timer

import (	
	"time"
)

var counter int = 15

func Decrement() {
	timer := time.NewTimer(time.Second)
	<-timer.C
	go func() {
	    counter--

	    if counter == 0 {
		/* Execute appropriate method here */	
	    }
	}()
}

func AddTime(amt int) {
	counter += amt
}










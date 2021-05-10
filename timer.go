package xzbase

import "time"

type TaskTimer interface {
	GetSecond() int64
	Run()
}

var timers = make([]TaskTimer, 0)

func AddTimer(t TaskTimer)  {
	timers = append(timers,t)
}

func StartTimer() {
	for _, t := range timers {

		go func(t TaskTimer) {
			tick := time.Tick(time.Second * time.Duration(t.GetSecond()))

			for {
				<-tick
				t.Run()
			}
		}(t)
	}
}

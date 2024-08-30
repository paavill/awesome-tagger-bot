package scheduler

import (
	"log"
	"testing"
	"time"
)

func TestCalcSleepTime(t *testing.T) {
	for nh := 0; nh < 24; nh++ {
		for nm := 0; nm < 60; nm++ {
			for sh := 0; sh < 24; sh++ {
				for sm := 0; sm < 60; sm++ {
					sleepTime := calcSleepTime(nh, nm, sh, sm)
					if sleepTime < 0 {
						t.Fail()
					} else {
						log.Printf("%02d:%02d -> %02d:%02d = %v", nh, nm, sh, sm, sleepTime)
						now := time.Duration(nh)*time.Hour + time.Duration(nm)*time.Minute + sleepTime
						setting := time.Duration(sh)*time.Hour + time.Duration(sm)*time.Minute

						delta := now - setting

						if delta.String() != "24h0m0s" && delta.String() != "0s" {
							log.Println(now.Hours(), now.Minutes(), setting.Hours(), setting.Minutes())
							t.Fail()
						}
						log.Println(now.String(), setting.String())
					}
				}
			}
		}
	}
}

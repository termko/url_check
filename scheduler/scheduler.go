package scheduler

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type Keks struct {
}

func (k Keks) Run() {
	fmt.Println("Wow every second...")
}

func Test() {
	c := cron.New()
	p := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	sch, _ := p.Parse("* * * * * *")
	c.Schedule(sch, &Keks{})
	c.Start()
	fmt.Println(c.Entries())
	time.Sleep(1000 * time.Second)
	fmt.Println("Wooooow")
}

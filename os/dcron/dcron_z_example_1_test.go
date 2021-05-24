package dcron_test

import (
	"github.com/osgochina/donkeygo/os/dcron"
	"github.com/osgochina/donkeygo/os/dlog"
	"time"
)

func Example_cronAddSingleton() {
	dcron.AddSingleton("* * * * * *", func() {
		dlog.Println("doing")
		time.Sleep(2 * time.Second)
	})
	select {}
}

package dcron_test

import (
	"github.com/osgochina/donkeygo/os/dcron"
	"time"

	"github.com/gogf/gf/os/glog"
)

func Example_cronAddSingleton() {
	dcron.AddSingleton("* * * * * *", func() {
		glog.Println("doing")
		time.Sleep(2 * time.Second)
	})
	select {}
}

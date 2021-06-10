package dcfg_test

import (
	"fmt"
	"github.com/osgochina/donkeygo/internal/intlog"
	"github.com/osgochina/donkeygo/os/dcfg"
)

func Example_mapSliceChange() {
	intlog.SetEnabled(false)
	defer intlog.SetEnabled(true)
	// For testing/example only.
	content := `{"map":{"key":"value"}, "slice":[59,90]}`
	dcfg.SetContent(content)
	defer dcfg.RemoveContent()

	m := dcfg.Instance().GetMap("map")
	fmt.Println(m)

	// Change the key-value pair.
	m["key"] = "john"

	// It changes the underlying key-value pair.
	fmt.Println(dcfg.Instance().GetMap("map"))

	s := dcfg.Instance().GetArray("slice")
	fmt.Println(s)

	// Change the value of specified index.
	s[0] = 100

	// It changes the underlying slice.
	fmt.Println(dcfg.Instance().GetArray("slice"))

	// output:
	// map[key:value]
	// map[key:john]
	// [59 90]
	// [100 90]
}

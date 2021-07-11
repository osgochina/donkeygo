package dhash_test

import (
	"github.com/osgochina/donkeygo/encoding/dhash"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

var (
	strBasic = []byte("This is the test string for hash.")
)

func Test_BKDRHash(t *testing.T) {
	var x uint32 = 200645773
	dtest.C(t, func(t *dtest.T) {
		j := dhash.BKDRHash(strBasic)
		t.Assert(j, x)
	})
}

func Test_BKDRHash64(t *testing.T) {
	var x uint64 = 4214762819217104013
	dtest.C(t, func(t *dtest.T) {
		j := dhash.BKDRHash64(strBasic)
		t.Assert(j, x)
	})
}

func Test_SDBMHash(t *testing.T) {
	var x uint32 = 1069170245
	dtest.C(t, func(t *dtest.T) {
		j := dhash.SDBMHash(strBasic)
		t.Assert(j, x)
	})
}

func Test_SDBMHash64(t *testing.T) {
	var x uint64 = 9881052176572890693
	dtest.C(t, func(t *dtest.T) {
		j := dhash.SDBMHash64(strBasic)
		t.Assert(j, x)
	})
}

func Test_RSHash(t *testing.T) {
	var x uint32 = 1944033799
	dtest.C(t, func(t *dtest.T) {
		j := dhash.RSHash(strBasic)
		t.Assert(j, x)
	})
}

func Test_RSHash64(t *testing.T) {
	var x uint64 = 13439708950444349959
	dtest.C(t, func(t *dtest.T) {
		j := dhash.RSHash64(strBasic)
		t.Assert(j, x)
	})
}

func Test_JSHash(t *testing.T) {
	var x uint32 = 498688898
	dtest.C(t, func(t *dtest.T) {
		j := dhash.JSHash(strBasic)
		t.Assert(j, x)
	})
}

func Test_JSHash64(t *testing.T) {
	var x uint64 = 13410163655098759877
	dtest.C(t, func(t *dtest.T) {
		j := dhash.JSHash64(strBasic)
		t.Assert(j, x)
	})
}

func Test_PJWHash(t *testing.T) {
	var x uint32 = 7244206
	dtest.C(t, func(t *dtest.T) {
		j := dhash.PJWHash(strBasic)
		t.Assert(j, x)
	})
}

func Test_PJWHash64(t *testing.T) {
	var x uint64 = 31150
	dtest.C(t, func(t *dtest.T) {
		j := dhash.PJWHash64(strBasic)
		t.Assert(j, x)
	})
}

func Test_ELFHash(t *testing.T) {
	var x uint32 = 7244206
	dtest.C(t, func(t *dtest.T) {
		j := dhash.ELFHash(strBasic)
		t.Assert(j, x)
	})
}

func Test_ELFHash64(t *testing.T) {
	var x uint64 = 31150
	dtest.C(t, func(t *dtest.T) {
		j := dhash.ELFHash64(strBasic)
		t.Assert(j, x)
	})
}

func Test_DJBHash(t *testing.T) {
	var x uint32 = 959862602
	dtest.C(t, func(t *dtest.T) {
		j := dhash.DJBHash(strBasic)
		t.Assert(j, x)
	})
}

func Test_DJBHash64(t *testing.T) {
	var x uint64 = 2519720351310960458
	dtest.C(t, func(t *dtest.T) {
		j := dhash.DJBHash64(strBasic)
		t.Assert(j, x)
	})
}

func Test_APHash(t *testing.T) {
	var x uint32 = 3998202516
	dtest.C(t, func(t *dtest.T) {
		j := dhash.APHash(strBasic)
		t.Assert(j, x)
	})
}

func Test_APHash64(t *testing.T) {
	var x uint64 = 2531023058543352243
	dtest.C(t, func(t *dtest.T) {
		j := dhash.APHash64(strBasic)
		t.Assert(j, x)
	})
}

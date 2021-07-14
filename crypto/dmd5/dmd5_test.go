package dmd5_test

import (
	"github.com/osgochina/donkeygo/crypto/dmd5"
	"github.com/osgochina/donkeygo/test/dtest"
	"os"
	"testing"
)

var (
	s = "pibigstar"
	// online generated MD5 value
	result = "d175a1ff66aedde64344785f7f7a3df8"
)

type user struct {
	name     string
	password string
	age      int
}

func TestEncrypt(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		encryptString, _ := dmd5.Encrypt(s)
		t.Assert(encryptString, result)

		result := "1427562bb29f88a1161590b76398ab72"
		encrypt, _ := dmd5.Encrypt(123456)
		t.AssertEQ(encrypt, result)
	})

	dtest.C(t, func(t *dtest.T) {
		user := &user{
			name:     "派大星",
			password: "123456",
			age:      23,
		}
		result := "70917ebce8bd2f78c736cda63870fb39"
		encrypt, _ := dmd5.Encrypt(user)
		t.AssertEQ(encrypt, result)
	})
}

func TestEncryptString(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		encryptString, _ := dmd5.EncryptString(s)
		t.Assert(encryptString, result)
	})
}

func TestEncryptFile(t *testing.T) {
	path := "test.text"
	errorPath := "err.txt"
	result := "e6e6e1cd41895beebff16d5452dfce12"
	dtest.C(t, func(t *dtest.T) {
		file, err := os.Create(path)
		defer os.Remove(path)
		defer file.Close()
		t.Assert(err, nil)
		_, _ = file.Write([]byte("Hello Go Frame"))
		encryptFile, _ := dmd5.EncryptFile(path)
		t.AssertEQ(encryptFile, result)
		// when the file is not exist,encrypt will return empty string
		errEncrypt, _ := dmd5.EncryptFile(errorPath)
		t.AssertEQ(errEncrypt, "")
	})

}

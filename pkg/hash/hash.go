package hash

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

func Verify(hpasswd string, passwd string, salt string) bool {
	fmt.Println("hpasswd: ", hpasswd, "passwd: ", passwd, "salt:", salt)
	cpasswd, _ := hex.DecodeString(hpasswd)
	csalt, _ := hex.DecodeString(salt)
	return Hash(passwd, string(csalt)) == string(cpasswd)
}
func Hash(passwd string, salt string) string {
	dk := pbkdf2.Key([]byte(passwd), []byte(salt), 1000, 64, sha1.New)
	return string(dk)
}

func Salt(l int) string {
	b := make([]byte, l)
	rand.Read(b)
	return string(b)
}

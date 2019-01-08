package GoRadius

import (
	"fmt"

	"golang.org/x/crypto/md4"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func NTHash(s string) string {
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	hasher := md4.New()
	t := transform.NewWriter(hasher, enc)
	t.Write([]byte(s))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

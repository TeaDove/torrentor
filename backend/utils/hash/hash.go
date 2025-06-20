package hash

import (
	"crypto/sha1"
	"encoding/base64"

	"github.com/pkg/errors"
)

func Sha1Base64Hash(s string) string {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(s))
	if err != nil {
		panic(errors.Wrap(err, "failed to write"))
	}

	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

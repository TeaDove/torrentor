package hash

import (
	"crypto/md5" //nolint: gosec // shoudn't be secure
	"encoding/base64"

	"github.com/pkg/errors"
)

func Sha1Base64Hash(s string) string {
	hasher := md5.New() //nolint: gosec // shoudn't be secure

	_, err := hasher.Write([]byte(s))
	if err != nil {
		panic(errors.Wrap(err, "failed to write"))
	}

	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

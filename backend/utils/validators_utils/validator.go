package validators_utils

import "github.com/go-playground/validator/v10"

var Validate = validator.New(validator.WithRequiredStructEnabled()) //nolint:gochecknoglobals // required for cache

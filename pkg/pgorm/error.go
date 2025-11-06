package pgorm

import "errors"

func IsRecordNotFound(err error) bool {
	return errors.Is(err, ErrRecordNotFound)
}

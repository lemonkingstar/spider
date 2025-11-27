package perror

import (
	"errors"
	"fmt"
)

// SerialWithError
// 依次执行所有函数，返回所有错误信息
func SerialWithError(fns ...func() error) func() error {
	return func() error {
		var errs error
		for _, fn := range fns {
			errs = AppendErr(errs, fn())
		}
		return errs
	}
}

// SerialUntilError
// 依次执行所有函数，报错立即返回
func SerialUntilError(fns ...func() error) func() error {
	return func() error {
		for _, fn := range fns {
			err := fn()
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func AppendErr(left, right error) error {
	if left == nil {
		return right
	} else if right == nil {
		return left
	}
	return fmt.Errorf("%w; %v", left, right)
}

func WrapErr(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	if len(args) == 0 {
		return fmt.Errorf("%w; %s", err, format)
	}
	return fmt.Errorf("%w; %s", err, fmt.Sprintf(format, args...))
}

func NewErr(msg string) error {
	return errors.New(msg)
}

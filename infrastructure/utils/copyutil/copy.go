package copyutil

import (
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

func Copy[T, R any](t *T) (*R, error) {
	r := new(R)
	if err := copier.Copy(r, t); err != nil {
		return nil, errors.Wrap(err, "copy error")
	}
	return r, nil
}

func CopySlice[T, R any](t []*T) ([]*R, error) {
	r := make([]*R, len(t))
	if err := copier.Copy(&r, &t); err != nil {
		return nil, errors.Wrap(err, "copy error")
	}
	return r, nil
}

func CopyTo[T, R any](from *T, to *R) error {
	if err := copier.Copy(to, from); err != nil {
		return errors.Wrap(err, "copy error")
	}
	return nil
}

func MustCopy[T, R any](t *T) *R {
	r, err := Copy[T, R](t)
	if err != nil {
		panic(err)
	}
	return r
}

func MustCopySlice[T, R any](t []*T) []*R {
	r, err := CopySlice[T, R](t)
	if err != nil {
		panic(err)
	}
	return r
}

func MustCopyTo[T, R any](from *T, to *R) {
	if err := CopyTo(from, to); err != nil {
		panic(err)
	}
}

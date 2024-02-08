package utils

import (
	"errors"

	"github.com/thoas/go-funk"
)

func Map[T, K any](input []T, predicate func(T) K) ([]K, error) {
	result, ok := funk.Map(input, predicate).([]K)
	if !ok {
		return nil, errors.New("error mapping")
	}

	return result, nil
}

func Contains[T any](input []T, predicate func(T) bool) bool {
	for _, elem := range input {
		if predicate(elem) {
			return true
		}
	}

	return false
}

package error

import (
	"errors"
	"fmt"
)

type ValidatorErrorBag struct {
	Errors map[string]string //A nicer way i could do this is make this multidimensional
}

func (v ValidatorErrorBag) Add(key, value string) {
	v.Errors[key] = value
}

func (v ValidatorErrorBag) Count() int {
	return len(v.Errors)
}

func (v ValidatorErrorBag) Get(key string) (string, error) {
	if v.Has(key) {
		return v.Errors[key], nil
	}

	formatted := fmt.Sprintf("Key, %s does not exist on this bag", key)

	return "", errors.New(formatted)
}

func (v ValidatorErrorBag) Has(key string) bool {

	_, exists := v.Errors[key]

	return exists
}

func (v ValidatorErrorBag) Reset() {
	for k, _ := range v.Errors {
		delete(v.Errors, k)
	}
}

func NewValidatorErrorBag() *ValidatorErrorBag {
	return &ValidatorErrorBag{Errors: make(map[string]string)}
}

package error

type ErrorBag interface {
	Add(key, value string)
	Delete(key string) bool
	Has(key string) bool
	Get(key string) (string, error)
	Reset() bool
	Count() int
}

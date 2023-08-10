package cfg

type Configurator[T any] func(el *T)

func Apply[T any](el *T, options []Configurator[T]) {
	for _, conf := range options {
		conf(el)
	}
}

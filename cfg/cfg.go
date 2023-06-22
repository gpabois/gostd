package cfg

import (
	"strconv"

	"github.com/gpabois/gostd/option"
)

type ConfigMap map[string]any

func (c ConfigMap) Provide() ConfigMap {
	return c
}

func GetInt(cfg *ConfigMap, path ...string) option.Option[int] {
	return option.Chain(Get(cfg, path...), func(val string) option.Option[int] {
		ival, err := strconv.Atoi(val)

		if err != nil {
			return option.None[int]()
		}

		return option.Some(ival)
	})
}

// Get the value from a Config Map
func Get(cfg *ConfigMap, path ...string) option.Option[string] {
	if len(path) == 0 {
		return option.None[string]()
	}

	if len(path) == 1 {
		val, ok := (*cfg)[path[0]]

		if !ok {
			return option.None[string]()
		}

		return option.Some(val.(string))
	}

	n, path := path[0], path[1:]

	val, ok := (*cfg)[n]

	if !ok {
		return option.None[string]()
	}

	node, ok := val.(ConfigMap)

	if !ok {
		return option.None[string]()
	}

	return Get(&node, path...)
}

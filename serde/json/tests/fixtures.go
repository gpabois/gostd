package json_tests

import "github.com/gpabois/gostd/option"

type simpleStruct struct {
	Opt option.Option[bool]
	Val int `serde:"val"`
}

package serde

import "fmt"

type UnhandledContentType struct {
	contentType string
}

func NewUnhandledContentType(contentType string) error {
	return UnhandledContentType{contentType: contentType}
}

func (err UnhandledContentType) Error() string {
	return fmt.Sprintf("unhandled content-type: %s", err.contentType)
}

type DeserializeError struct {
	internal error
}

func NewDeserializeError(err error) error {
	return DeserializeError{internal: err}
}

func (err DeserializeError) Error() string {
	return fmt.Sprintf("deserialize error %s: ", err.internal)
}

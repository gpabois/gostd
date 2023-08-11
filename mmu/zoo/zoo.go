package zoo

import "github.com/gpabois/gostd/option"

// An objet allocator, but the resource is managed by the allocator itself.
// The resource can be accessed.
type Zoo[ID any, T any, PT interface{ *T }] interface {
	// Try to allocate a new element, returns None if it cannot.
	TryNew() option.Option[ID]
	// Allocate a new element, panic if it cannot.
	New() ID
	// Delete the element
	Delete(ID)
	// Allows to access the resource, the pointer must not exit the scope
	With(id ID, fn func(ptr PT))
	// Check if the resource exists
	Exists(id ID) bool
}

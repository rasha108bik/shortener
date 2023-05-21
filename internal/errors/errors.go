// error package contains response errors with a description
package errors

import "errors"

// Generic api handlers errors.
var (
	ErrOriginalURLExist = errors.New("originalURL exist") // url is exist
	ErrNoSuchID         = errors.New("no such id")        // no such id
	ErrURLDeleted       = errors.New("URL deleted")       // url deleted
)

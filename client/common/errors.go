package common

import "errors"

var (
	ErrCacheEntryExpired = errors.New("entry is expired")
	ErrCannotLoadPosts   = errors.New("cannot load posts")
	ErrNotFound          = errors.New("not found")
	// ErrNotAuthenticated means the server bounced us to a login page or
	// answered 401, so no usable session was presented.
	ErrNotAuthenticated = errors.New("not authenticated")
	// ErrForbidden means the server answered 403. Reddit returns this for a
	// rejected session cookie ("Blocked"), but also for private or banned
	// content, so the two cannot be distinguished from the status alone.
	ErrForbidden             = errors.New("forbidden")
	ErrCannotOpenCacheFile   = errors.New("cannot open cache file")
	ErrCannotEncodeCacheFile = errors.New("cannot encode cache file")
	ErrCannotDecodeCacheFile = errors.New("cannot decode cache file")
)

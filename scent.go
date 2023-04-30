package scent

import (
	"fmt"
	"sync"
)

var (
	ErrLoadFailed    = fmt.Errorf("failed to load a scene")
	ErrAlreadyLoaded = fmt.Errorf("%w: already loaded", ErrLoadFailed)
	ErrUnloadFailed  = fmt.Errorf("failed to unload a scene")
	ErrNotLoaded     = fmt.Errorf("%w: was not loaded", ErrUnloadFailed)
)

type Scene[L, U, D any] interface {
	Load(L) (unload func() error, err error)
	Update(U) error
	Draw(D)
}

type Switch[L, U, D any] struct {
	stack []loaded[L, U, D]
	mu    sync.Mutex
}

// Updates all scenes top-to-bottom
func (s *Switch[L, U, D]) Update(v U) error {
}

// Draws all scenes bottom-to-top
func (s *Switch[L, U, D]) Draw(v D) {
func (s *Switch[L, U, D]) LoadScene(v L, scene Scene[L, U, D]) error {


}
}

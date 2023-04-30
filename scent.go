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

type loaded[L, U, D any] struct {
	Scene[L, U, D]
	Unload func() error
}

type Switch[L, U, D any] struct {
	stack []loaded[L, U, D]
	mu    sync.Mutex
}

// Updates all scenes top-to-bottom
func (s *Switch[L, U, D]) Update(v U) error {
	stack := s.stack
	for i := 0; i < len(stack); i++ {
		err := stack[i].Update(v)
		if err != nil {
			return err
		}
	}

	return nil
}

// Draws all scenes bottom-to-top
func (s *Switch[L, U, D]) Draw(v D) {
	stack := s.stack
	for i := len(stack); i > 0; i-- {
		stack[i-1].Draw(v)
	}
}

func (s *Switch[L, U, D]) LoadScene(v L, scene Scene[L, U, D]) error {
	s.mu.Lock()
	for _, l := range s.stack {
		if l.Scene == scene {
			s.mu.Unlock()
			return ErrAlreadyLoaded
		}
	}
	s.mu.Unlock()

	unload, err := scene.Load(v)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrLoadFailed, err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// realloc to allow removed scenes to get garbage collected
	new := make([]loaded[L, U, D], len(s.stack)+1)
	new[0] = loaded[L, U, D]{scene, unload}
	copy(new[1:], s.stack)
	s.stack = new

	return nil
}

func (s *Switch[L, U, D]) UnloadScene(v L, scene Scene[L, U, D]) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var i int
	var unload func() error
	for j, s := range s.stack {
		if s.Scene == scene {
			i = j
			unload = s.Unload
			goto found
		}
	}
	return ErrNotLoaded

found:

	if unload != nil {
		s.mu.Unlock()
		err := unload()
		s.mu.Lock()

		if err != nil {
			return fmt.Errorf("%w: %w", ErrUnloadFailed, err)
		}
	}

	// realloc to allow removed scenes to get garbage collected
	new := make([]loaded[L, U, D], len(s.stack)-1)
	copy(new[:i], s.stack[:i])
	copy(new[i:], s.stack[i+1:])
	s.stack = new

	return nil
}

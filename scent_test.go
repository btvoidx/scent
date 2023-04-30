package scent_test

import (
	"testing"

	"github.com/btvoidx/scent"
)

type none struct{}

type Scene struct {
	load   func() error
	unload func() error
	update func() error
	draw   func()
}

func (s *Scene) Load(none) (func() error, error) {
	if err := s.load(); err != nil {
		return nil, err
	}

	return s.unload, nil
}

func (s *Scene) Update(none) error {
	if s.update != nil {
		return s.update()
	}
	return nil
}

func (s *Scene) Draw(none) {
	if s.draw != nil {
		s.draw()
	}
}

var _ scent.Scene[none, none, none] = (*Scene)(nil)

func TestLoad(t *testing.T) {
	s := new(scent.Switch[none, none, none])

	var loaded, unloaded bool
	scene := &Scene{
		load:   func() error { loaded = true; return nil },
		unload: func() error { unloaded = true; return nil },
	}

	// Loading
	if err := s.LoadScene(none{}, scene); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if err := s.LoadScene(none{}, scene); err == nil {
		t.Fatalf("expected err; got nil")
	}

	// Unloading
	if err := s.UnloadScene(none{}, scene); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if err := s.UnloadScene(none{}, scene); err == nil {
		t.Fatalf("expected err; got nil")
	}

	switch {
	case !loaded:
		t.Fatalf("was not loaded")
	case !unloaded:
		t.Fatalf("was not unloaded")
	}
}

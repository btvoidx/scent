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
	if s.load != nil {
		return s.unload, s.load()
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

func TestUpdateOrder(t *testing.T) {
	s := new(scent.Switch[none, none, none])

	var res int
	s.LoadScene(none{}, &Scene{update: func() error { res *= res; return nil }})
	s.LoadScene(none{}, &Scene{update: func() error { res *= 3; return nil }})
	s.LoadScene(none{}, &Scene{update: func() error { res += 1; return nil }})

	s.Update(none{})

	if res != 9 {
		t.Fatalf("wrong update order")
	}
}

func TestDrawOrder(t *testing.T) {
	s := new(scent.Switch[none, none, none])

	var res int
	s.LoadScene(none{}, &Scene{draw: func() { res += 1 }})
	s.LoadScene(none{}, &Scene{draw: func() { res *= 3 }})
	s.LoadScene(none{}, &Scene{draw: func() { res *= res }})

	s.Draw(none{})

	if res != 9 {
		t.Fatalf("wrong draw order")
	}
}

package scent

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	ErrLoadFailed    = fmt.Errorf("failed to load a scene")
	ErrAlreadyLoaded = fmt.Errorf("%w: already loaded", ErrLoadFailed)
	ErrUnloadFailed  = fmt.Errorf("failed to unload a scene")
	ErrNotLoaded     = fmt.Errorf("%w: was not loaded", ErrUnloadFailed)
)

type Scene[G any] interface {
	Load(G) (unload func() error, err error)
	Update() error
	Draw(*ebiten.Image)
}

type Switch[G any] struct {
	stack []Scene[G]
}

// Updates all scenes top-to-bottom
func (Switch[G]) Update() error

// Draws all scenes bottom-to-top
func (Switch[G]) Draw(*ebiten.Image)

func (Switch[G]) LoadScene(Scene[G]) error

func (Switch[G]) UnloadScene(Scene[G]) error

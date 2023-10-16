package main

import (
	"fmt"
	"os"
)

type operation interface {
	Apply() error
	String() string
	AppendSuffix(suffix string) error
}

type rename struct {
	src string
	dst string
}

func (r *rename) Apply() error   { return os.Rename(r.src, r.dst) }
func (r *rename) String() string { return fmt.Sprintf("Rename '%s' to '%s'", r.src, r.dst) }
func (r *rename) AppendSuffix(suffix string) error {
	err := os.Rename(r.src, r.src+suffix)
	if err != nil {
		return fmt.Errorf("failed to add temporary suffix to '%s': %w", r.src, err)
	}
	r.src += suffix
	return nil
}

type remove struct{ src string }

func (r *remove) Apply() error   { return os.RemoveAll(r.src) }
func (r *remove) String() string { return fmt.Sprintf("Remove '%s'", r.src) }
func (r *remove) AppendSuffix(suffix string) error {
	err := os.Rename(r.src, r.src+suffix)
	if err != nil {
		return fmt.Errorf("failed to add temporary suffix to '%s': %w", r.src, err)
	}
	r.src += suffix
	return nil
}

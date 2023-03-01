package main

import (
	"fmt"
	"os"
)

type operation interface {
	Apply() error
	String() string
}

type rename struct {
	src string
	dst string
}

func (r *rename) Apply() error   { return os.Rename(r.src, r.dst) }
func (r *rename) String() string { return fmt.Sprintf("Rename '%s' to '%s'", r.src, r.dst) }

type remove struct{ src string }

func (r *remove) Apply() error   { return os.RemoveAll(r.src) }
func (r *remove) String() string { return fmt.Sprintf("Remove '%s'", r.src) }

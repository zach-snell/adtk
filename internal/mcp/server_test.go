package mcp

import (
	"testing"

	"github.com/zach-snell/adtk/internal/devops"
)

func TestNewServer_RegistersTools(t *testing.T) {
	t.Parallel()
	c := devops.NewTestClient("http://localhost:0")
	s := newServer(c)
	if s == nil {
		t.Fatal("newServer returned nil")
	}
}

func TestNewServer_WithRealConstructor(t *testing.T) {
	t.Parallel()
	s := New("test-org", "test-pat")
	if s == nil {
		t.Fatal("New returned nil")
	}
}

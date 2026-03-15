package test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()

	// Exit with the appropriate exit code
	os.Exit(exitVal)
}

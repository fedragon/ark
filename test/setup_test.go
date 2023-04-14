package test

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("Setting up test environment...")

	code := m.Run()

	fmt.Println("Tearing down test environment...")

	_ = os.Remove("ark.db")
	_ = os.RemoveAll("archive")

	os.Exit(code)
}

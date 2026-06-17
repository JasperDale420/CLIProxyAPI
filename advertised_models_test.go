package main

import (
	"os"
	"strings"
	"testing"
)

func TestZAIAdvertisesGLM52(t *testing.T) {
	for _, path := range []string{"README.md", "README_CN.md", "README_JA.md", "config.example.yaml"} {
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(strings.ToLower(string(body)), "glm-5.2") {
			t.Fatalf("%s does not advertise glm-5.2", path)
		}
	}
}

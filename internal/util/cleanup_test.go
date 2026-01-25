package util

import "testing"

func TestRemoveAiLinks(t *testing.T) {
	input := "Start ([link](https://example.com)) middle ([another](https://example.org)) end"
	expected := "Start  middle  end"
	result := RemoveAiLinks(input)
	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRemoveAiLinksNoMatch(t *testing.T) {
	input := "Nothing to change here"
	result := RemoveAiLinks(input)
	if result != input {
		t.Fatalf("expected %q, got %q", input, result)
	}
}

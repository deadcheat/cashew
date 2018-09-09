package strings

import "testing"

func TestFirstString(t *testing.T) {
	testee := []string{"hoge", "fuga"}
	if FirstString(testee) != "hoge" {
		t.Error("didn't return first string of slice")
	}

	if FirstString(nil) != "" {
		t.Error("didn't return empty string")
	}

	if FirstString([]string{}) != "" {
		t.Error("didn't return empty string")
	}
}

func TestStringSliceContainsTrue(t *testing.T) {
	// true pattern
	testee := []string{
		"false",
		"hoge",
		"true",
	}
	if !StringSliceContainsTrue(testee) {
		t.Error("it should return true")
	}

	// false pattern
	testee = []string{
		"100",
	}
	if StringSliceContainsTrue(testee) {
		t.Error("it should return false")
	}
}

package main

import "testing"

func TestIsMongoConnectionString(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"mongodb://localhost:27017", true},
		{"mongodb+srv://atlas.test", true},
		{"postgres://localhost:5432", false},
	}

	for _, c := range cases {
		got := isMongoConnectionString(c.input)
		if got != c.want {
			t.Errorf("isMongoConnectionString(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

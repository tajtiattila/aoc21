package main

import "testing"

func TestDay10(t *testing.T) {
	tests := []struct {
		src  string
		want int
	}{
		{"[({(<(())[]>[[{[]{<()<>>", 288957},
		{"[(()[<>])]({[<{<<[]>>(", 5566},
		{"(((({<>}<{<{<>}{[]{[]{}", 1480781},
		{"{<[[]]>}<{[{[{[]{()[[[]", 995444},
		{"<{([{{}}[<[[[<>{}]]]>[]]", 294},
	}

	for _, tt := range tests {
		_, got := bracer(tt.src)
		if got != tt.want {
			t.Fatalf("%q: %d != %d", tt.src, got, tt.want)
		}
	}
}

package main

import "testing"

func TestDay11a(t *testing.T) {
	o := parsedumboct(5, `11111
19991
19191
19991
11111`)

	o.step()
	o.step()

	want := `45654
51115
61116
51115
45654
`
	if want != o.String() {
		t.Fatalf("got: \n%s\n\nwant: %s\n", o.String(), want)
	}

}

package git

import "testing"

func TestShortenSHA(t *testing.T) {
	commitID := "6a6bcddc365ca3a38c9055a603c9590a7fae7ca6"

	wanted := "6a6bcd"
	if s := ShortenSHA(commitID); s != wanted {
		t.Fatalf("ShortenSHA got %s, wanted %s", s, wanted)
	}
}

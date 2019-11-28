package git

// ShortenSHA trims a SHA to the first 6 characters.
func ShortenSHA(s string) string {
	return s[:6]
}

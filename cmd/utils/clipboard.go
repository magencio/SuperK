package utils

// Clipboard represents the clipboard to copy and paste strings across the app
type Clipboard struct {
	Content string
}

// NewClipboard creates a new Clipboard
func NewClipboard() *Clipboard {
	return &Clipboard{}
}

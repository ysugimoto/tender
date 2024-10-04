package tender

type RenderOption func(t *Template)

func WithHtmlEscape() RenderOption {
	return func(t *Template) {
		t.enableEscape = true
	}
}

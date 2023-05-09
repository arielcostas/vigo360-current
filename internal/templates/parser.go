package templates

import (
	gmf "github.com/mdigger/goldmark-figures"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var parser goldmark.Markdown = goldmark.New(
	goldmark.WithExtensions(extension.Footnote),
	goldmark.WithExtensions(extension.Typographer),
	goldmark.WithExtensions(gmf.Extension),
)

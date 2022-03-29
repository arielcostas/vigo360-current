/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package common

import (
	goldmarkfigures "github.com/mdigger/goldmark-figures"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var Parser goldmark.Markdown = goldmark.New(
	goldmark.WithExtensions(extension.Footnote),
	goldmark.WithExtensions(goldmarkfigures.Extension),
)

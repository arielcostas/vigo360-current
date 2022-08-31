// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package internal

import (
	"fmt"
	"net/http"
)

func (s *Server) handlePublicIndexnowKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s\n", r.URL.Path[1:len(r.URL.Path)-4])
	}
}

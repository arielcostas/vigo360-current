/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import (
	"encoding/json"
	"net/http"

	"vigo360.es/new/internal/logger"
)

type jsonError struct {
	Error    bool   `json:"error"`
	ErrorMsg string `json:"errorMsg"`
}

type jsonHandler func(http.ResponseWriter, *http.Request) *appError

func (fn jsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	if ae := fn(w, r); ae != nil {
		var rid = r.Context().Value("rid").(string)

		logger.Error("[%s] `%s` %s: %s", rid, r.URL.Path, ae.Message, ae.Error.Error())

		bytes, err := json.MarshalIndent(jsonError{true, ae.Response}, "", "\t")
		if err != nil {
			logger.Error("error producing error JSON: %s", err.Error())
			return
		}

		w.Header().Add("Vigo360-RID", rid)
		w.WriteHeader(ae.Status)
		w.Write(bytes)
	}
}

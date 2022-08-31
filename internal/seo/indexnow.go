// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package seo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type indexnowRequestBody struct {
	Host    string   `json:"host"`
	Key     string   `json:"key"`
	UrlList []string `json:"urlList"`
}

func BingIndexnowRequest(urls []string) error {
	var client = http.DefaultClient
	requestBytes, err := json.Marshal(indexnowRequestBody{
		Host:    strings.TrimPrefix(os.Getenv("DOMAIN"), "https://"),
		Key:     os.Getenv("INDEXNOW_KEY"),
		UrlList: urls,
	})
	if err != nil {
		return err
	}
	var requestBody = bytes.NewBuffer(requestBytes)

	request, err := http.NewRequest(http.MethodPost, "https://www.bing.com/indexnow", requestBody)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		return err
	}
	return nil
}

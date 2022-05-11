/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import "os"

func fullCanonica(path string) string {
	return os.Getenv("DOMAIN") + path
}

func getMinimo(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

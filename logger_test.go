// SPDX-License-Identifier: MIT

package logs

var (
	_ Logger = &Entry{}
	_ Logger = &logger{}
	_ Logger = &emptyLogger{}
)

// SPDX-License-Identifier: MIT

package logs

var (
	_ Logger = &entry{}
	_ Logger = &logger{}
	_ Logger = &emptyLogger{}
)

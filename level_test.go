// SPDX-License-Identifier: MIT

package logs

import (
	"encoding"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/issue9/assert/v3"
)

var (
	testLevel = LevelInfo

	_ encoding.TextMarshaler   = LevelDebug
	_ encoding.TextUnmarshaler = &testLevel
	_ fmt.Stringer             = LevelError
)

func TestIsValidLevel(t *testing.T) {
	a := assert.New(t, false)

	a.False(IsValidLevel(-1)).
		False(IsValidLevel(LevelFatal + 1)).
		True(IsValidLevel(LevelError))
}

func TestParseLevel(t *testing.T) {
	a := assert.New(t, false)

	lv, err := ParseLevel("INFO")
	a.NotError(err).Equal(lv, LevelInfo)

	lv, err = ParseLevel("erro")
	a.NotError(err).Equal(lv, LevelError)

	lv, err = ParseLevel("not-exists")
	a.ErrorString(err, "无效的值").Equal(lv, -1)

	lv, err = ParseLevel(levelStrings[-1])
	a.ErrorString(err, "无效的值").Equal(lv, -1)
}

func TestLevel_UnmarshalText(t *testing.T) {
	a := assert.New(t, false)

	l := LevelWarn
	a.NotError(json.Unmarshal([]byte(`"inFo"`), &l))
	a.Equal(l, LevelInfo)

	l = LevelWarn
	a.ErrorString(json.Unmarshal([]byte(`"not-exists"`), &l), "无效的值")
}

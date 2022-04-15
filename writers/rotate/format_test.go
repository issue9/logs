// SPDX-License-Identifier: MIT

package rotate

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/issue9/assert/v2"
)

func TestCutString(t *testing.T) {
	a := assert.New(t, false)

	p, s, err := cutString("")
	a.Empty(p).Empty(s).Equal(err, ErrIndexNotExists)

	p, s, err = cutString("%i")
	a.NotError(err).
		Empty(p).
		Empty(s)

	p, s, err = cutString("test%i")
	a.NotError(err).
		Equal(p, "test").
		Empty(s)

	p, s, err = cutString("test-200602%i06y01-15")
	a.NotError(err).
		Equal(p, "test-200602").
		Equal(s, "06y01-15")

	p, s, err = cutString("test-200602/06y01-15")
	a.ErrorString(err, "不能包含")
}

func TestGetIndex(t *testing.T) {
	a := assert.New(t, false)
	now := time.Now()
	prefixValue := now.Format("2006.")
	suffixValue := now.Format(".01")

	i, err := getIndex("./testdata", prefixValue, suffixValue)
	a.NotError(err).Equal(i, 0)

	w := func(i int) {
		name := "./testdata/" + prefixValue + strconv.Itoa(i) + suffixValue
		a.NotError(os.WriteFile(name, []byte("123"), os.ModePerm))
	}

	w(5)
	i, err = getIndex("./testdata", prefixValue, suffixValue)
	a.NotError(err).Equal(i, 5)

	w(8)
	i, err = getIndex("./testdata", prefixValue, suffixValue)
	a.NotError(err).Equal(i, 8)
}

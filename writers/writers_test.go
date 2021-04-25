// SPDX-License-Identifier: MIT

package writers

import "bytes"

type testContainer struct {
	bytes.Buffer
	closed, flushed bool
}

func (c *testContainer) Close() error {
	c.flushed = true
	c.closed = true
	return nil
}

func (c *testContainer) Flush() error {
	c.flushed = true
	return nil
}

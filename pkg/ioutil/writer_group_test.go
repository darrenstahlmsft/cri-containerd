/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ioutil

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type writeCloser struct {
	buf    bytes.Buffer
	closed bool
}

func (wc *writeCloser) Write(p []byte) (int, error) {
	return wc.buf.Write(p)
}

func (wc *writeCloser) Close() error {
	wc.closed = true
	return nil
}

func TestEmptyWriterGroup(t *testing.T) {
	wg := NewWriterGroup()
	_, err := wg.Write([]byte("test"))
	assert.Error(t, err)
}

func TestClosedWriterGroup(t *testing.T) {
	wg := NewWriterGroup()
	wc := &writeCloser{}
	key, data := "test key", "test data"

	err := wg.Add(key, wc)
	assert.NoError(t, err)

	n, err := wg.Write([]byte(data))
	assert.Equal(t, len(data), n)
	assert.Equal(t, data, wc.buf.String())
	assert.NoError(t, err)

	wg.Close()
	assert.True(t, wc.closed)

	err = wg.Add(key, &writeCloser{})
	assert.Error(t, err)

	_, err = wg.Write([]byte(data))
	assert.Error(t, err)
}

func TestAddRemoveWriter(t *testing.T) {
	wg := NewWriterGroup()
	wc1, wc2 := &writeCloser{}, &writeCloser{}
	key1, key2 := "test key 1", "test key 2"

	err := wg.Add(key1, wc1)
	assert.NoError(t, err)
	_, err = wg.Write([]byte("test data 1"))
	assert.NoError(t, err)
	assert.Equal(t, "test data 1", wc1.buf.String())

	err = wg.Add(key2, wc2)
	assert.NoError(t, err)
	_, err = wg.Write([]byte("test data 2"))
	assert.NoError(t, err)
	assert.Equal(t, "test data 1test data 2", wc1.buf.String())
	assert.Equal(t, "test data 2", wc2.buf.String())

	wg.Remove(key1)
	_, err = wg.Write([]byte("test data 3"))
	assert.NoError(t, err)
	assert.Equal(t, "test data 1test data 2", wc1.buf.String())
	assert.Equal(t, "test data 2test data 3", wc2.buf.String())

	wg.Close()
}

// Copyright (C) Kumo inc. and its affiliates.
// Author: Jeff.li lijippy@163.com
// All rights reserved.
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package gzip

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kumose-go/archive/config"
	"github.com/stretchr/testify/require"
)

func TestGzFile(t *testing.T) {
	tmp := t.TempDir()
	f, err := os.Create(filepath.Join(tmp, "test.gz"))
	require.NoError(t, err)
	defer f.Close()
	archive := New(f)
	defer archive.Close()

	require.NoError(t, archive.Add(config.File{
		Destination: "sub1/sub2/subfoo.txt",
		Source:      "../testdata/sub1/sub2/subfoo.txt",
	}))
	require.EqualError(t, archive.Add(config.File{
		Destination: "foo.txt",
		Source:      "../testdata/foo.txt",
	}), "gzip: failed to add foo.txt, only one file can be archived in gz format")
	require.NoError(t, archive.Close())
	require.NoError(t, f.Close())

	f, err = os.Open(f.Name())
	require.NoError(t, err)
	defer f.Close()

	info, err := f.Stat()
	require.NoError(t, err)
	require.Lessf(t, info.Size(), int64(500), "archived file should be smaller than %d", info.Size())

	gzf, err := gzip.NewReader(f)
	require.NoError(t, err)
	defer gzf.Close()

	require.Equal(t, "sub1/sub2/subfoo.txt", gzf.Name)

	bts, err := io.ReadAll(gzf)
	require.NoError(t, err)
	require.Equal(t, "sub\n", string(bts))
}

func TestGzFileCustomMtime(t *testing.T) {
	f, err := os.Create(filepath.Join(t.TempDir(), "test.gz"))
	require.NoError(t, err)
	defer f.Close()
	archive := New(f)
	defer archive.Close()

	now := time.Now().Truncate(time.Second)

	require.NoError(t, archive.Add(config.File{
		Destination: "sub1/sub2/subfoo.txt",
		Source:      "../testdata/sub1/sub2/subfoo.txt",
		Info: config.FileInfo{
			ParsedMTime: now,
		},
	}))
	require.NoError(t, archive.Close())
	require.NoError(t, f.Close())

	f, err = os.Open(f.Name())
	require.NoError(t, err)
	defer f.Close()

	info, err := f.Stat()
	require.NoError(t, err)
	require.Lessf(t, info.Size(), int64(500), "archived file should be smaller than %d", info.Size())

	gzf, err := gzip.NewReader(f)
	require.NoError(t, err)
	defer gzf.Close()

	require.Equal(t, "sub1/sub2/subfoo.txt", gzf.Name)
	require.Equal(t, now, gzf.ModTime)
}

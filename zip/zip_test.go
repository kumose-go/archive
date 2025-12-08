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

package zip

import (
	"archive/zip"
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kumose-go/archive/config"
	"github.com/kumose-go/archive/testlib"
	"github.com/stretchr/testify/require"
)

func TestZipFile(t *testing.T) {
	tmp := t.TempDir()
	f, err := os.Create(filepath.Join(tmp, "test.zip"))
	require.NoError(t, err)
	defer f.Close()
	archive := New(f)
	defer archive.Close()

	require.Error(t, archive.Add(config.File{
		Source:      "../testdata/nope.txt",
		Destination: "nope.txt",
	}))
	require.NoError(t, archive.Add(config.File{
		Source:      "../testdata/foo.txt",
		Destination: "foo.txt",
	}))
	require.NoError(t, archive.Add(config.File{
		Source:      "../testdata/sub1",
		Destination: "sub1",
	}))
	require.NoError(t, archive.Add(config.File{
		Source:      "../testdata/sub1/bar.txt",
		Destination: "sub1/bar.txt",
	}))
	require.NoError(t, archive.Add(config.File{
		Source:      "../testdata/sub1/executable",
		Destination: "sub1/executable",
	}))
	require.NoError(t, archive.Add(config.File{
		Source:      "../testdata/sub1/sub2",
		Destination: "sub1/sub2",
	}))
	require.NoError(t, archive.Add(config.File{
		Source:      "../testdata/sub1/sub2/subfoo.txt",
		Destination: "sub1/sub2/subfoo.txt",
	}))
	require.NoError(t, archive.Add(config.File{
		Source:      "../testdata/regular.txt",
		Destination: "regular.txt",
	}))
	require.NoError(t, archive.Add(config.File{
		Source:      "../testdata/link.txt",
		Destination: "link.txt",
	}))

	require.ErrorIs(t, archive.Add(config.File{
		Source:      "../testdata/regular.txt",
		Destination: "link.txt",
	}), fs.ErrExist)

	require.NoError(t, archive.Close())
	require.Error(t, archive.Add(config.File{
		Source:      "tar.go",
		Destination: "tar.go",
	}))
	require.NoError(t, f.Close())

	f, err = os.Open(f.Name())
	require.NoError(t, err)
	defer f.Close()

	info, err := f.Stat()
	require.NoError(t, err)
	require.Lessf(t, info.Size(), int64(1000), "archived file should be smaller than %d", info.Size())

	r, err := zip.NewReader(f, info.Size())
	require.NoError(t, err)

	paths := make([]string, len(r.File))
	for i, zf := range r.File {
		paths[i] = zf.Name
		if zf.Name == "sub1/executable" && !testlib.IsWindows() {
			require.NotEqualf(
				t,
				0,
				zf.Mode()&0o111,
				"expected executable perms, got %s",
				zf.Mode().String(),
			)
		}
		if zf.Name == "link.txt" {
			require.NotEqual(t, 0, zf.FileInfo().Mode()&os.ModeSymlink)
			rc, err := zf.Open()
			require.NoError(t, err)
			var link bytes.Buffer
			_, err = io.Copy(&link, rc)
			require.NoError(t, err)
			rc.Close()
			require.Equal(t, "regular.txt", link.String())
		}
	}
	require.Equal(t, []string{
		"foo.txt",
		"sub1/bar.txt",
		"sub1/executable",
		"sub1/sub2/subfoo.txt",
		"regular.txt",
		"link.txt",
	}, paths)
}

func TestZipFileInfo(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	f, err := os.Create(filepath.Join(t.TempDir(), "test.zip"))
	require.NoError(t, err)
	defer f.Close()
	archive := New(f)
	defer archive.Close()

	require.NoError(t, archive.Add(config.File{
		Source:      "../testdata/foo.txt",
		Destination: "nope.txt",
		Info: config.FileInfo{
			Mode:        0o755,
			Owner:       "carlos",
			Group:       "root",
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

	r, err := zip.NewReader(f, info.Size())
	require.NoError(t, err)

	require.Len(t, r.File, 1)
	for _, next := range r.File {
		require.Equal(t, "nope.txt", next.Name)
		require.Equal(t, now.Unix(), next.Modified.Unix())
		require.Equal(t, fs.FileMode(0o755), next.FileInfo().Mode())
	}
}

func TestTarInvalidLink(t *testing.T) {
	archive := New(io.Discard)
	defer archive.Close()

	require.NoError(t, archive.Add(config.File{
		Source:      "../testdata/badlink.txt",
		Destination: "badlink.txt",
	}))
}

// TODO: add copying test

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

package archive

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/kumose-go/archive/config"
	"github.com/kumose-go/archive/testlib"
	"github.com/stretchr/testify/require"
)

func TestArchive(t *testing.T) {
	folder := t.TempDir()
	empty, err := os.Create(folder + "/empty.txt")
	require.NoError(t, err)
	require.NoError(t, empty.Close())
	require.NoError(t, os.Mkdir(folder+"/folder-inside", 0o755))

	for _, format := range []string{"tar.gz", "zip", "gz", "tar.xz", "tar", "tgz", "txz", "tar.zst", "tzst"} {
		t.Run(format, func(t *testing.T) {
			f1, err := os.Create(filepath.Join(t.TempDir(), "1.tar"))
			require.NoError(t, err)

			archive, err := New(f1, format)
			require.NoError(t, err)
			require.NoError(t, archive.Add(config.File{
				Source:      empty.Name(),
				Destination: "empty.txt",
			}))
			require.Error(t, archive.Add(config.File{
				Source:      empty.Name() + "_nope",
				Destination: "dont.txt",
			}))
			require.NoError(t, archive.Close())
			require.NoError(t, f1.Close())

			if format == "tar.xz" || format == "txz" || format == "gz" || format == "tar.zst" || format == "tzst" {
				_, err := Copy(f1, io.Discard, format)
				require.Error(t, err)
				return
			}

			f1, err = os.Open(f1.Name())
			require.NoError(t, err)
			f2, err := os.Create(filepath.Join(t.TempDir(), "2.tar"))
			require.NoError(t, err)

			a, err := Copy(f1, f2, format)
			require.NoError(t, err)
			require.NoError(t, f1.Close())

			require.NoError(t, a.Add(config.File{
				Source:      empty.Name(),
				Destination: "added_later.txt",
			}))
			require.NoError(t, a.Add(config.File{
				Source:      empty.Name(),
				Destination: "ملف.txt",
			}))
			require.NoError(t, a.Close())
			require.NoError(t, f2.Close())

			require.ElementsMatch(
				t,
				[]string{"empty.txt", "added_later.txt", "ملف.txt"},
				testlib.LsArchive(t, f2.Name(), format),
			)
		})
	}

	// unsupported format...
	t.Run("7z", func(t *testing.T) {
		_, err := New(io.Discard, "7z")
		require.EqualError(t, err, "invalid archive format: 7z")
	})
}

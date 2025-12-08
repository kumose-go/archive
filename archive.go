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
	"fmt"
	"io"
	"os"

	"github.com/kumose-go/archive/config"
	"github.com/kumose-go/archive/gzip"
	"github.com/kumose-go/archive/tar"
	"github.com/kumose-go/archive/targz"
	"github.com/kumose-go/archive/tarxz"
	"github.com/kumose-go/archive/tarzst"
	"github.com/kumose-go/archive/zip"
)

// Archive represents a compression archive files from disk can be written to.
type Archive interface {
	Close() error
	Add(f config.File) error
}

// New archive.
func New(w io.Writer, format string) (Archive, error) {
	switch format {
	case "tar.gz", "tgz":
		return targz.New(w), nil
	case "tar":
		return tar.New(w), nil
	case "gz":
		return gzip.New(w), nil
	case "tar.xz", "txz":
		return tarxz.New(w), nil
	case "tar.zst", "tzst":
		return tarzst.New(w), nil
	case "zip":
		return zip.New(w), nil
	}
	return nil, fmt.Errorf("invalid archive format: %s", format)
}

// Copy copies the source archive into a new one, which can be appended at.
// Source needs to be in the specified format.
func Copy(r *os.File, w io.Writer, format string) (Archive, error) {
	switch format {
	case "tar.gz", "tgz":
		return targz.Copy(r, w)
	case "tar":
		return tar.Copy(r, w)
	case "zip":
		return zip.Copy(r, w)
	}
	return nil, fmt.Errorf("invalid archive format: %s", format)
}

// Unarchive extracts the source archive to destination according to format.
// opts controls strip top dir and overwrite behavior.
func Unarchive(src, dest, format string, opts config.ExtractOptions) error {
	switch format {
	case "tar.gz", "tgz":
		return targz.ExtractTargz(src, dest, opts)
	case "tar":
		return tar.ExtractTar(src, dest, opts)
	case "tar.xz", "txz":
		return tarxz.ExtractTarXZ(src, dest, opts)
	case "tar.zst", "tzst":
		return tarzst.ExtractTarZST(src, dest, opts)
	case "zip":
		return zip.ExtractZip(src, dest, opts)
	default:
		return fmt.Errorf("unsupported archive format: %s", format)
	}
}

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

package targz

import (
	"compress/gzip"
	"io"

	"github.com/kumose-go/archive/config"
	"github.com/kumose-go/archive/tar"
)

// Archive as tar.gz.
type Archive struct {
	gw *gzip.Writer
	tw *tar.Archive
}

// New tar.gz archive.
func New(target io.Writer) Archive {
	// the error will be nil since the compression level is valid
	gw, _ := gzip.NewWriterLevel(target, gzip.BestCompression)
	tw := tar.New(gw)
	return Archive{
		gw: gw,
		tw: &tw,
	}
}

func Copy(source io.Reader, target io.Writer) (Archive, error) {
	// the error will be nil since the compression level is valid
	gw, _ := gzip.NewWriterLevel(target, gzip.BestCompression)
	srcgz, err := gzip.NewReader(source)
	if err != nil {
		return Archive{}, err
	}
	tw, err := tar.Copy(srcgz, gw)
	return Archive{
		gw: gw,
		tw: &tw,
	}, err
}

// Close all closeables.
func (a Archive) Close() error {
	if err := a.tw.Close(); err != nil {
		return err
	}
	return a.gw.Close()
}

// Add file to the archive.
func (a Archive) Add(f config.File) error {
	return a.tw.Add(f)
}

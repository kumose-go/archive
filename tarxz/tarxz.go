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

package tarxz

import (
	"io"

	"github.com/kumose-go/archive/config"
	"github.com/kumose-go/archive/tar"
	"github.com/ulikunitz/xz"
)

// Archive as tar.xz.
type Archive struct {
	xzw *xz.Writer
	tw  *tar.Archive
}

// New tar.xz archive.
func New(target io.Writer) Archive {
	xzw, _ := xz.WriterConfig{DictCap: 16 * 1024 * 1024}.NewWriter(target)
	tw := tar.New(xzw)
	return Archive{
		xzw: xzw,
		tw:  &tw,
	}
}

// Close all closeables.
func (a Archive) Close() error {
	if err := a.tw.Close(); err != nil {
		return err
	}
	return a.xzw.Close()
}

// Add file to the archive.
func (a Archive) Add(f config.File) error {
	return a.tw.Add(f)
}

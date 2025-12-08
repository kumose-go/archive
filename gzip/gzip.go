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
	"fmt"
	"io"
	"os"

	gzip "github.com/klauspost/pgzip"
	"github.com/kumose-go/archive/config"
)

// Archive as gz.
type Archive struct {
	gw *gzip.Writer
}

// New gz archive.
func New(target io.Writer) Archive {
	// the error will be nil since the compression level is valid
	gw, _ := gzip.NewWriterLevel(target, gzip.BestCompression)
	return Archive{
		gw: gw,
	}
}

// Close all closeables.
func (a Archive) Close() error {
	return a.gw.Close()
}

// Add file to the archive.
func (a Archive) Add(f config.File) error {
	if a.gw.Name != "" {
		return fmt.Errorf("gzip: failed to add %s, only one file can be archived in gz format", f.Destination)
	}
	file, err := os.Open(f.Source) // #nosec
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	a.gw.Name = f.Destination
	if f.Info.ParsedMTime.IsZero() {
		a.gw.ModTime = info.ModTime()
	} else {
		a.gw.ModTime = f.Info.ParsedMTime
	}
	_, err = io.Copy(a.gw, file)
	return err
}

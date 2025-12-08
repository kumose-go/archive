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

package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/kumose-go/archive/config"
)

// Archive as tar.
type Archive struct {
	tw    *tar.Writer
	files map[string]bool
}

// New tar archive.
func New(target io.Writer) Archive {
	return Archive{
		tw:    tar.NewWriter(target),
		files: map[string]bool{},
	}
}

// Copy creates a new tar with the contents of the given tar.
func Copy(source io.Reader, target io.Writer) (Archive, error) {
	w := New(target)
	r := tar.NewReader(source)
	for {
		header, err := r.Next()
		if err == io.EOF || header == nil {
			break
		}
		if err != nil {
			return Archive{}, err
		}
		w.files[header.Name] = true
		if err := w.tw.WriteHeader(header); err != nil {
			return w, err
		}
		if _, err := io.Copy(w.tw, r); err != nil {
			return w, err
		}
	}
	return w, nil
}

// Close all closeables.
func (a Archive) Close() error {
	return a.tw.Close()
}

// Add file to the archive.
func (a Archive) Add(f config.File) error {
	if _, ok := a.files[f.Destination]; ok {
		return &fs.PathError{Err: fs.ErrExist, Path: f.Destination, Op: "add"}
	}
	a.files[f.Destination] = true
	info, err := os.Lstat(f.Source) // #nosec
	if err != nil {
		return fmt.Errorf("%s: %w", f.Source, err)
	}
	var link string
	if info.Mode()&os.ModeSymlink != 0 {
		link, err = os.Readlink(f.Source) // #nosec
		if err != nil {
			return fmt.Errorf("%s: %w", f.Source, err)
		}
	}
	header, err := tar.FileInfoHeader(info, link)
	if err != nil {
		return fmt.Errorf("%s: %w", f.Source, err)
	}
	header.Name = f.Destination
	if !f.Info.ParsedMTime.IsZero() {
		header.ModTime = f.Info.ParsedMTime
	}
	if f.Info.Mode != 0 {
		header.Mode = int64(f.Info.Mode)
	}
	if f.Info.Owner != "" {
		header.Uid = 0
		header.Uname = f.Info.Owner
	}
	if f.Info.Group != "" {
		header.Gid = 0
		header.Gname = f.Info.Group
	}
	if err = a.tw.WriteHeader(header); err != nil {
		return fmt.Errorf("%s: %w", f.Source, err)
	}
	if info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil
	}
	file, err := os.Open(f.Source) // #nosec
	if err != nil {
		return fmt.Errorf("%s: %w", f.Source, err)
	}
	defer file.Close()
	if _, err := io.Copy(a.tw, file); err != nil {
		return fmt.Errorf("%s: %w", f.Source, err)
	}
	return nil
}

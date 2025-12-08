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
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kumose-go/archive/config"
)

// ExtractTargz extracts a .tar.gz archive to the destination directory.
// Supports options to strip the top-level directory and overwrite existing files.
func ExtractTargz(src, dest string, opts config.ExtractOptions) error {
	// Open tar.gz file
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	var topDir string
	first := true

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		name := header.Name
		// Determine top-level dir for StripTopDir
		if first && opts.StripTopDir {
			topDir = strings.SplitN(name, string(os.PathSeparator), 2)[0]
			first = false
		}

		if opts.StripTopDir {
			name = strings.TrimPrefix(name, topDir+"/")
		}
		targetPath := filepath.Join(dest, name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if !opts.Overwrite {
				if _, err := os.Lstat(targetPath); err == nil {
					return fmt.Errorf("file exists: %s", targetPath)
				}
			}
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}
			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		case tar.TypeSymlink:
			if !opts.Overwrite {
				if _, err := os.Lstat(targetPath); err == nil {
					return fmt.Errorf("symlink exists: %s", targetPath)
				}
			}
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}
			if err := os.Symlink(header.Linkname, targetPath); err != nil {
				return err
			}
		default:
			// skip other types
			continue
		}
	}

	return nil
}

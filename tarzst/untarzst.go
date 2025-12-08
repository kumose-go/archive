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

package tarzst

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
	"github.com/kumose-go/archive/config"
)

// ExtractTarZST extracts a .tar.zst archive to dest directory.
func ExtractTarZST(src, dest string, opts config.ExtractOptions) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	dec, err := zstd.NewReader(file)
	if err != nil {
		return err
	}
	defer dec.Close()

	tr := tar.NewReader(dec)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, header.Name)

		// Strip top-level directory if requested
		if opts.StripTopDir {
			parts := strings.SplitN(header.Name, string(os.PathSeparator), 2)
			if len(parts) == 2 {
				target = filepath.Join(dest, parts[1])
			} else {
				// top-level file/dir, skip if nothing left
				continue
			}
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if !opts.Overwrite {
				if _, err := os.Stat(target); err == nil {
					return fmt.Errorf("file %s exists", target)
				}
			} else {
				_ = os.MkdirAll(filepath.Dir(target), 0755)
			}
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
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
				if _, err := os.Lstat(target); err == nil {
					return fmt.Errorf("symlink %s exists", target)
				}
			}
			_ = os.MkdirAll(filepath.Dir(target), 0755)
			if err := os.Symlink(header.Linkname, target); err != nil {
				return err
			}
		default:
			// skip other types for simplicity
		}
	}

	return nil
}

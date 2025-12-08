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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kumose-go/archive/config"
)

// ExtractZip extracts a .zip archive to dest directory.
func ExtractZip(src, dest string, opts config.ExtractOptions) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		name := f.Name

		// Strip top-level directory if requested
		if opts.StripTopDir {
			parts := strings.SplitN(name, string(os.PathSeparator), 2)
			if len(parts) == 2 {
				name = parts[1]
			} else {
				// top-level file, skip if nothing left
				continue
			}
		}

		target := filepath.Join(dest, name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, f.Mode()); err != nil {
				return err
			}
			continue
		}

		// create parent directories
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		if !opts.Overwrite {
			if _, err := os.Stat(target); err == nil {
				return fmt.Errorf("file %s exists", target)
			}
		}

		outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		if _, err := io.Copy(outFile, rc); err != nil {
			rc.Close()
			outFile.Close()
			return err
		}

		rc.Close()
		outFile.Close()
	}

	return nil
}

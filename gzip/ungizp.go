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
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kumose-go/archive/config"
)

// ExtractGzip extracts a .gz file to the destination directory.
// If StripTopDir is true, the output file will be stripped of its path (only filename used).
// Overwrite controls whether existing files are replaced.
func ExtractGzip(src, dest string, opts config.ExtractOptions) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	// Determine output file path
	outName := gzReader.Name
	if outName == "" {
		// fallback to source filename without .gz
		base := filepath.Base(src)
		if len(base) > 3 && base[len(base)-3:] == ".gz" {
			outName = base[:len(base)-3]
		} else {
			outName = base + ".out"
		}
	}

	if opts.StripTopDir {
		outName = filepath.Base(outName)
	}

	targetPath := filepath.Join(dest, outName)

	if !opts.Overwrite {
		if _, err := os.Lstat(targetPath); err == nil {
			return fmt.Errorf("file exists: %s", targetPath)
		}
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, gzReader); err != nil {
		return err
	}

	return nil
}

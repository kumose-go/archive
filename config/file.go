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

package config

import (
	"os"
	"time"
)

// File is a file inside an archive.
type File struct {
	Source      string   `yaml:"src,omitempty" json:"src,omitempty"`
	Destination string   `yaml:"dst,omitempty" json:"dst,omitempty"`
	StripParent bool     `yaml:"strip_parent,omitempty" json:"strip_parent,omitempty"`
	Info        FileInfo `yaml:"info,omitempty" json:"info,omitempty"`
	Default     bool     `yaml:"-" json:"-"`
}

// FileInfo is the file info of a file.
type FileInfo struct {
	Owner       string      `yaml:"owner,omitempty" json:"owner,omitempty"`
	Group       string      `yaml:"group,omitempty" json:"group,omitempty"`
	Mode        os.FileMode `yaml:"mode,omitempty" json:"mode,omitempty"`
	MTime       string      `yaml:"mtime,omitempty" json:"mtime,omitempty"`
	ParsedMTime time.Time   `yaml:"-" json:"-"`
}

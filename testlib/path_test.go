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

package testlib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckPath(t *testing.T) {
	requireSkipped := func(tb testing.TB, skipped bool) {
		tb.Helper()
		t.Cleanup(func() {
			require.Equalf(tb, skipped, tb.Skipped(), "expected skipped to be %v", skipped)
		})
	}

	t.Run("local", func(t *testing.T) {
		t.Setenv("CI", "false")

		t.Run("in path", func(t *testing.T) {
			requireSkipped(t, false)
			if IsWindows() {
				CheckPath(t, "cmd.exe")
			} else {
				CheckPath(t, "echo")
			}
		})

		t.Run("not in path", func(t *testing.T) {
			requireSkipped(t, true)
			CheckPath(t, "do-not-exist")
		})
	})

	t.Run("CI", func(t *testing.T) {
		t.Setenv("CI", "true")

		t.Run("in path on CI", func(t *testing.T) {
			requireSkipped(t, false)
			CheckPath(t, "echo")
		})

		t.Run("not in path on CI", func(t *testing.T) {
			requireSkipped(t, false)
			CheckPath(t, "do-not-exist")
		})
	})
}

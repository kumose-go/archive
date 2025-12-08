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
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
)

// CheckPath skips the test if the binary is not in the PATH, or if CI is true.
func CheckPath(tb testing.TB, cmd string) {
	tb.Helper()
	if !InPath(cmd) {
		tb.Skipf("%s not in PATH", cmd)
	}
}

// IsCI returns true if we have the "CI" environment variable set to true.
func IsCI() bool {
	return os.Getenv("CI") == "true"
}

// InPath returns true if the given cmd is in the PATH, or if CI is true.
func InPath(cmd string) bool {
	if IsCI() {
		return true
	}
	_, err := exec.LookPath(cmd)
	return err == nil
}

// IsWindows returns true if current OS is Windows.
func IsWindows() bool { return runtime.GOOS == "windows" }

// SkipIfWindows skips the test if runtime OS is windows.
func SkipIfWindows(tb testing.TB, args ...any) {
	tb.Helper()
	if IsWindows() {
		tb.Skip(args...)
	}
}

// Echo returns a `echo s` command, handling it on windows.
func Echo(s string) string {
	if IsWindows() {
		return "cmd.exe /c echo " + s
	}
	return "echo " + s
}

// Touch returns a `touch name` command, handling it on windows.
func Touch(name string) string {
	if IsWindows() {
		return "cmd.exe /c copy nul " + name
	}
	return "touch " + name
}

// ShC returns the command line for the given cmd wrapped into a `sh -c` in
// linux/mac, and the cmd.exe command in windows.
func ShC(cmd string) string {
	if IsWindows() {
		return fmt.Sprintf("cmd.exe /c '%s'", cmd)
	}
	return fmt.Sprintf("sh -c '%s'", cmd)
}

// Exit returns a command that exits the given status, handling windows.
func Exit(status int) string {
	if IsWindows() {
		return fmt.Sprintf("cmd.exe /c exit /b %d", status)
	}
	return fmt.Sprintf("exit %d", status)
}

// Copyright 2021 smzgl@foxmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"syscall"

	"github.com/coreos/go-semver/semver"
)

// SortSemanticVersion is used to sort semantic version
func SortSemanticVersion(items []string) ([]string, []string) {
	versionMap := make(map[*semver.Version]string, len(items))
	versions := make(semver.Versions, 0, len(items))
	var malformed []string
	for _, item := range items {
		s := item
		if strings.HasPrefix(s, "v") {
			s = item[1:]
		}
		version, err := semver.NewVersion(s)
		if err != nil {
			malformed = append(malformed, item)
			continue
		}
		versionMap[version] = item
		versions = append(versions, version)
	}
	sort.Sort(versions)
	var data []string
	for _, version := range versions {
		data = append(data, versionMap[version])
	}
	sort.Strings(malformed)
	return malformed, data
}

// DeduplicateSliceStably is used to deduplicate slice items stably
func DeduplicateSliceStably(items []string) []string {
	data := make([]string, 0, len(items))
	deduplicate := map[string]struct{}{}
	for _, val := range items {
		if _, exists := deduplicate[val]; !exists {
			deduplicate[val] = struct{}{}
			data = append(data, val)
		}
	}
	return data
}

// ContainsEmpty is used to check whether items contains empty string
func ContainsEmpty(items ...string) bool {
	return Contains(items, "")
}

// Contains is used to check whether the target is in items
func Contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

// SetToSlice is used to convert set<string> to slice<string>
func SetToSlice(set map[string]struct{}) []string {
	data := make([]string, 0, len(set))
	for key := range set {
		data = append(data, key)
	}
	return data
}

// GetMapKeys is used to get the keys of map
func GetMapKeys(dict map[string]string) []string {
	data := make([]string, 0, len(dict))
	for key := range dict {
		data = append(data, key)
	}
	return data
}

// GetExitCode is used to parse exit code from cmd error
func GetExitCode(err error) int {
	if exitErr, ok := err.(*exec.ExitError); ok {
		// The program has exited with an exit code != 0
		// This works on both Unix and Windows. Although package
		// syscall is generally platform dependent, WaitStatus is
		// defined for both Unix and Windows and in both cases has
		// an ExitStatus() method with the same signature.
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 1
}

var (
	regexpEnvironmentVar = regexp.MustCompile(`\$[A-Za-z_]+`)
	regexpRegularVersion = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`)
)


// IsRegularVersion is used to determine whether the version number is a regular version number
// Regular: va.b.c, and a, b, c are all numbers
func IsRegularVersion(s string) bool {
	return regexpRegularVersion.MatchString(s)
}

// RenderWithEnv is used to render string with env
func RenderWithEnv(s string, ext map[string]string) string {
	matches := regexpEnvironmentVar.FindAllString(s, -1)
	for _, match := range matches {
		key := match[1:]
		val := ext[key]
		if val == "" {
			val = os.Getenv(key)
		}
		if val != "" {
			s = strings.ReplaceAll(s, match, val)
		}
	}
	return s
}

// RenderPathWithEnv is used to render path with environment
func RenderPathWithEnv(path string, ext map[string]string) string {
	return filepath.Clean(RenderWithEnv(path, ext))
}

// SplitGoPackageVersion is used to split go package version
func SplitGoPackageVersion(pkg string) (path string, version string, ok bool) {
	i := strings.Index(pkg, "@")
	if i == -1 {
		return "", "", false
	}
	return pkg[:i], pkg[i+1:], true
}

// JoinGoPackageVersion is used to join go path and versions
func JoinGoPackageVersion(path, version string) string {
	return strings.Join([]string{
		path, version,
	}, "@")
}

// GetBinaryFileName is used to get os based binary file name
func GetBinaryFileName(name string) string {
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(name, ".exe") {
			return name + ".exe"
		}
		return name
	}
	return name
}

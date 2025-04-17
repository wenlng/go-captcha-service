package helper

import "path"

// GetResourceDirAbsPath 。
func GetResourceDirAbsPath() string {
	return path.Join(GetPWD(), "resources")
}

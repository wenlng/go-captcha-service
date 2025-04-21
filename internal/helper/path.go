/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package helper

import "path"

// GetResourceDirAbsPath 。
func GetResourceDirAbsPath() string {
	return path.Join(GetPWD(), "resources")
}

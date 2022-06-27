package file_utils

import "strings"

func GetSuffixPaths(path string, suffix string, excludePaths ...string) []string {
	var r []string
	if IsDir(path) && !checkExcludePaths(path, excludePaths...) {
		if strings.HasSuffix(path, suffix) {
			r = append(r, path)
		}

		s := ListDir(path)
		for _, v := range s {
			t := GetSuffixPaths(v, suffix)
			r = append(r, t...)
		}
	} else {
		if strings.HasSuffix(path, suffix) {
			r = append(r, path)
		}
	}

	return r
}

func checkExcludePaths(path string, excludePaths ...string) bool {
	for _, v := range excludePaths {
		if path == v {
			return true
		}
	}
	return false
}
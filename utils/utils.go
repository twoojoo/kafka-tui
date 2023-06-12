package utils

import "strconv"

func BytesToString(bytes int) string {
	m := 1024 //multiplier

	if bytes < 1024 {
		return strconv.Itoa(bytes) + " bytes"
	} else if bytes < m*m {
		return strconv.Itoa(bytes/m) + " KB"
	} else if bytes < m*m*m {
		return strconv.Itoa((bytes/m)/m) + " MB"
	} else if bytes < m*m*m*m {
		return strconv.Itoa(((bytes/m)/m)/m) + " GB"
	} else {
		return strconv.Itoa((((bytes/m)/m)/m)/m) + " TB"
	}
}
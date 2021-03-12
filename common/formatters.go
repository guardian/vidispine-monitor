package common

import "fmt"

func FormatBytes(byteValue int64) string {
	suffixes := []string{"bytes", "Kib", "Mib", "Gib", "Tib", "Eib"}

	reducedValue := byteValue
	for _, suffix := range suffixes {
		if reducedValue < 1024 {
			return fmt.Sprintf("%d%s", reducedValue, suffix)
		} else {
			reducedValue = reducedValue / 1024
		}
	}
	return fmt.Sprintf("%dEib", reducedValue)
}

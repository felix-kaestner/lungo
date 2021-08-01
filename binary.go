package lungo

import "fmt"

func byteToBinaryIEC(b int64) string {
	const base = 1024
	if b < base {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(base), 0
	for n := b / base; n >= base && exp < 5; n /= base {
		div *= base
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

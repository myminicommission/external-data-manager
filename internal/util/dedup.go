package util

import (
	"github.com/myminicommission/external-data-manager/internal/games"
)

// RemoveDuplicateMinis removes... duplicate... minis... ¯\_(ツ)_/¯ (based on Mini.Name)
func RemoveDuplicateMinis(minis []games.Mini) []games.Mini {
	keys := make(map[string]bool)
	list := []games.Mini{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range minis {
		if _, value := keys[entry.Name]; !value {
			keys[entry.Name] = true
			list = append(list, entry)
		}
	}
	return list
}

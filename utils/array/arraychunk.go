package array

import "hotel-engine/utils"

func Chunks(ids []string, size int) [][]string {
	var batch [][]string
	for i := 0; i < len(ids); i += size {
		batch = append(batch, ids[i:utils.Min(i+size, len(ids))])
	}
	return batch
}

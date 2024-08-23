package util

func SplitDataIntoBatch[T any](data []T, batchSize int) [][]T {
	if data == nil || len(data) < 0 {
		return make([][]T, 0)
	}

	dataBatches := make([][]T, 0)

	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}

		dataBatches = append(dataBatches, data[i:end])
	}

	return dataBatches
}

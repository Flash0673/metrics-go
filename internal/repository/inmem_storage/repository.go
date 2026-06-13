package inmem_storage

type MemStorage struct {
	storageGauge   map[string]float64
	storageCounter map[string]int64
}

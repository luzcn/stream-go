package stream

type PartitionOffsets map[string]string

type offsetStore interface {
	SetOffset(int32, int64) error
	GetOffsets() (*PartitionOffsets, error)
}

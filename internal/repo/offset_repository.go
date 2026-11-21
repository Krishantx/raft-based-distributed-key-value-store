package repo

type Offset_repository struct {
	offset map[uint64]uint64
}

func New_Offset_Repository() *Offset_repository {

	o := Offset_repository{
		offset: make(map[uint64]uint64),
	}

	return &o
}

func (o *Offset_repository) AppendOffset(index uint64, offset uint64) {
	o.offset[index] = offset
}

func (o *Offset_repository) GetOffset(index uint64) uint64 {
	return o.offset[index]
}

func (o *Offset_repository) TotalItems() int {
	return len(o.offset)
}

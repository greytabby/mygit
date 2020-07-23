package objects

type Blob struct {
	data []byte
}

func NewBlob(data []byte) *Blob {
	return &Blob{data}
}

func (o *Blob) Type() []byte {
	return []byte("blob")
}

func (o *Blob) Serialize() []byte {
	return o.data
}

func (o *Blob) Deserialize(data []byte) {
	o.data = data
}

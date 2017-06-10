package storage

import (
	"encoding/binary"
	"errors"
)

// ErrInvalidInfoData is the error returned when invalid info data provided.
var ErrInvalidInfoData = errors.New("storage: invalid info data")

// Info torrent info struct
type Info struct {
	Hash       []byte // 20 bytes
	Incomplete uint32
	Complete   uint32
	Downloaded uint32
	Name       []byte
}

// Marshal marshals a torrent Info to a byte slice
func (i *Info) Marshal() []byte {
	data := make([]byte, 32+len(i.Name))
	copy(data, i.Hash[:20])
	binary.LittleEndian.PutUint32(data[20:], i.Incomplete)
	binary.LittleEndian.PutUint32(data[24:], i.Complete)
	binary.LittleEndian.PutUint32(data[28:], i.Downloaded)
	copy(data[32:], i.Name)

	return data
}

// Unmarshal unmarshals data to a torrent Info.
func (i *Info) Unmarshal(data []byte) error {
	if len(data) < 32 {
		return ErrInvalidInfoData
	}

	i.Hash = data[:20]
	i.Incomplete = binary.LittleEndian.Uint32(data[20:])
	i.Complete = binary.LittleEndian.Uint32(data[24:])
	i.Downloaded = binary.LittleEndian.Uint32(data[28:])
	i.Name = data[32:]

	return nil
}

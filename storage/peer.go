package storage

import (
	"encoding/binary"
	"time"
)

type Peer struct {
	ID         []byte // 20 bytes
	Uploaded   uint64
	Downloaded uint64
	Left       uint64
	IP         []byte // 4 bytes
	Port       uint16
	Key        []byte
	LastUpdate time.Time
}

func (p *Peer) Marshal() []byte {
	data := make([]byte, 58+len(p.Key))
	copy(data, p.ID[:20])
	binary.LittleEndian.PutUint64(data[20:], p.Uploaded)
	binary.LittleEndian.PutUint64(data[28:], p.Downloaded)
	binary.LittleEndian.PutUint64(data[36:], p.Left)
	copy(data[44:], p.IP[:4])
	binary.LittleEndian.PutUint16(data[48:], p.Port)
	binary.LittleEndian.PutUint64(data[50:], uint64(p.LastUpdate.UnixNano()))
	copy(data[58:], p.Key)

	return data
}

func (p *Peer) Unmarshal(data []byte) error {
	if len(data) < 58 {
		return ErrInvalidInfoData
	}

	p.ID = data[:20]
	p.Uploaded = binary.LittleEndian.Uint64(data[20:])
	p.Downloaded = binary.LittleEndian.Uint64(data[28:])
	p.Left = binary.LittleEndian.Uint64(data[36:])
	p.IP = data[44:48]
	p.Port = binary.LittleEndian.Uint16(data[48:])
	p.LastUpdate = time.Unix(0, int64(binary.LittleEndian.Uint64(data[50:])))
	p.Key = data[58:]

	return nil
}

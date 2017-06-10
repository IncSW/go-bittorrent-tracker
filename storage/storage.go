package storage

import "errors"

var (
	ErrInvalidInfoData = errors.New("storage: invalid info data")
	ErrTorrentNotFound = errors.New("storage: torrent not found")
)

type Storage interface {
	GetInfo([]byte) (*Info, error)
	SetInfo(*Info) error
	SetPeer([]byte, *Peer) error
	GetActivePeers([]byte) ([]*Peer, error)
}

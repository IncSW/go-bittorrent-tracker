package storage

import "errors"

// ErrTorrentNotFound is the error returned when torrent not found.
var ErrTorrentNotFound = errors.New("storage: torrent not found")

// Storage interface
type Storage interface {
	// GetInfo returns torrent info by infoHash.
	GetInfo([]byte) (*Info, error)

	// SetInfo store torrent info.
	SetInfo(*Info) error

	// SetPeer store peer identified by infoHash.
	SetPeer([]byte, *Peer) error

	// GetActivePeers returns active peers for torrent by infoHash.
	GetActivePeers([]byte) ([]*Peer, error)
}

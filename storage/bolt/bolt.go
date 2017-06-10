package bolt

import (
	"bytes"
	"errors"
	"time"

	"github.com/IncSW/go-bittorrent-tracker/storage"
	"github.com/boltdb/bolt"
)

var (
	boltInfoKey = []byte("_info")

	errBoltFound = errors.New("found")
)

type boltSotrage struct {
	db *bolt.DB
}

func (s *boltSotrage) GetInfo(hash []byte) (*storage.Info, error) {
	info := &storage.Info{}
	if err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(hash)
		if bucket == nil {
			return storage.ErrTorrentNotFound
		}

		return info.Unmarshal(bucket.Get(boltInfoKey))
	}); err != nil {
		return nil, err
	}

	return info, nil
}

func (s *boltSotrage) SetInfo(info *storage.Info) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(info.Hash)
		if err != nil {
			return err
		}

		return bucket.Put(boltInfoKey, info.Marshal())
	})
}

func (s *boltSotrage) SetPeer(infoHash []byte, peer *storage.Peer) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(infoHash)
		if bucket == nil {
			return storage.ErrTorrentNotFound
		}

		return bucket.Put(peer.ID, peer.Marshal())
	})
}

func (s *boltSotrage) GetActivePeers(infoHash []byte) ([]*storage.Peer, error) {
	var peers []*storage.Peer
	now := time.Now()

	if err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(infoHash)
		if bucket == nil {
			return storage.ErrTorrentNotFound
		}

		if err := bucket.ForEach(func(key []byte, value []byte) error {
			if bytes.Equal(key, boltInfoKey) {
				return nil
			}

			peer := &storage.Peer{}
			if err := peer.Unmarshal(value); err != nil {
				return err
			}

			if now.Sub(peer.LastUpdate) > time.Hour {
				return nil
			}

			peers = append(peers, peer)

			return nil
		}); err != nil && err != errBoltFound {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return peers, nil
}

func New(db *bolt.DB) storage.Storage {
	return &boltSotrage{
		db: db,
	}
}

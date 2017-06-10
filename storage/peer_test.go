package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPeerEncoding(t *testing.T) {
	assert := assert.New(t)

	peer := &Peer{
		ID:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		Uploaded:   10000,
		Downloaded: 100000,
		Left:       100,
		IP:         []byte{1, 2, 3, 4},
		Port:       12345,
		Key:        []byte("key"),
		LastUpdate: time.Now(),
	}
	data := peer.Marshal()

	peerUnmarshaled := &Peer{}
	err := peerUnmarshaled.Unmarshal(data)
	if !assert.NoError(err) ||
		!assert.Equal(peer, peerUnmarshaled) ||
		!assert.Equal(data, peerUnmarshaled.Marshal()) {
		return
	}
}

func BenchmarkPeerMarshal(b *testing.B) {
	peer := &Peer{
		ID:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		Uploaded:   10000,
		Downloaded: 100000,
		Left:       100,
		IP:         []byte{1, 2, 3, 4},
		Port:       12345,
		Key:        []byte("key"),
		LastUpdate: time.Now(),
	}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		peer.Marshal()
	}
}

func BenchmarkPeerUnmarshal(b *testing.B) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 16, 39, 0, 0, 0, 0, 0, 0, 160, 134, 1, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 57, 48, 216, 172, 83, 132, 224, 152, 196, 20, 107, 101, 121}
	peer := &Peer{}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		peer.Unmarshal(data)
	}
}

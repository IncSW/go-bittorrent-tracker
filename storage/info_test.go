package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfoEncoding(t *testing.T) {
	assert := assert.New(t)

	info := &Info{
		Hash:       []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		Incomplete: 1000,
		Complete:   10000,
		Downloaded: 100000,
		Name:       []byte("Torrent name"),
	}
	data := info.Marshal()

	infoUnmarshaled := &Info{}
	err := infoUnmarshaled.Unmarshal(data)
	if !assert.NoError(err) ||
		!assert.Equal(info, infoUnmarshaled) ||
		!assert.Equal(data, infoUnmarshaled.Marshal()) {
		return
	}
}

func BenchmarkInfoMarshal(b *testing.B) {
	info := &Info{
		Hash:       []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		Incomplete: 1000,
		Complete:   10000,
		Downloaded: 100000,
		Name:       []byte("Torrent name"),
	}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		info.Marshal()
	}
}

func BenchmarkInfoUnmarshal(b *testing.B) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 232, 3, 0, 0, 16, 39, 0, 0, 160, 134, 1, 0, 84, 111, 114, 114, 101, 110, 116, 32, 110, 97, 109, 101}
	info := &Info{}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		info.Unmarshal(data)
	}
}

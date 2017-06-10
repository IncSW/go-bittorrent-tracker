package http

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"net"
	"time"

	bencode "github.com/IncSW/go-bencode"
	"github.com/IncSW/go-bittorrent-tracker/storage"
	"github.com/valyala/fasthttp"
)

func (s *httpServer) announce(ctx *fasthttp.RequestCtx) {
	request, err := makeRequest(ctx)
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	info, err := s.storage.GetInfo(request.infoHash)
	if err != nil {
		if err != storage.ErrTorrentNotFound {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
			return
		}

		info = &storage.Info{
			Hash: request.infoHash,
		}
	}

	switch request.event {
	case eventStarted:
		if request.left == 0 {
			info.Complete++
		} else {
			info.Incomplete++
		}
	case eventCompleted:
		info.Complete++
		if info.Incomplete != 0 {
			info.Incomplete--
		}
		info.Downloaded++
	case eventStopped:
		if request.left == 0 {
			if info.Complete != 0 {
				info.Complete--
			}
		} else {
			if info.Incomplete != 0 {
				info.Incomplete--
			}
		}
	}

	if err = s.storage.SetInfo(info); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	if err = s.storage.SetPeer(info.Hash, &storage.Peer{
		ID:         request.peerID,
		Uploaded:   request.uploaded,
		Downloaded: request.downloaded,
		Left:       request.left,
		IP:         request.ip,
		Port:       request.port,
		Key:        request.key,
		LastUpdate: time.Now(),
	}); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	peers, err := s.storage.GetActivePeers(info.Hash)
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	var randomPeers []*storage.Peer
	totalPeers := uint8(0)
	rand.Seed(time.Now().UnixNano())
	for _, i := range rand.Perm(len(peers)) {
		if bytes.Equal(request.peerID, peers[i].ID) {
			continue
		}

		randomPeers = append(randomPeers, peers[i])
		totalPeers++
		if totalPeers >= request.numWant {
			break
		}
	}

	data := map[string]interface{}{
		"interval":     600,
		"min interval": 300,
		// "tracker id": "wtf?",
		"complete":   info.Complete,
		"incomplete": info.Incomplete,
	}

	if totalPeers != 0 {
		if request.compact {
			offset := 0
			peersBencode := make([]byte, totalPeers*6)
			for _, peer := range randomPeers {
				copy(peersBencode[offset:], peer.IP)
				offset += 4
				binary.BigEndian.PutUint16(peersBencode[offset:], peer.Port)
				offset += 2
			}
			data["peers"] = peersBencode
		} else {
			peersBencode := []interface{}{}
			for _, peer := range randomPeers {
				peerBencode := map[string]interface{}{
					"ip":   net.IP(peer.IP).String(),
					"port": peer.Port,
				}
				if !request.noPeerID {
					peerBencode["peer id"] = peer.ID
				}
				peersBencode = append(peersBencode, peerBencode)
			}
			data["peers"] = peersBencode
		}
	}

	body, err := bencode.Marshal(data)
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("text/plain; charset=utf8")
	ctx.Write(body)
}

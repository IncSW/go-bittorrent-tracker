package http

import (
	"bytes"
	"net"

	bencode "github.com/IncSW/go-bencode"
	"github.com/IncSW/go-bittorrent-tracker/protocol"
	"github.com/IncSW/go-bittorrent-tracker/storage"
	"github.com/valyala/fasthttp"
)

var (
	infoHashKey = []byte("info_hash")

	announcePath = []byte("/announce")
	scrapePath   = []byte("/scrape")

	failureInfoHashNotFound = []byte("info hash not found")
	failureInvalidInfoHash  = []byte("invalid info hash")
	failureTorrentsNotFound = []byte("torrents not found")
)

type httpServer struct {
	storage storage.Storage
	server  *fasthttp.Server
}

func (s *httpServer) ListenAndServe(address string) error {
	tcpAddress, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp4", tcpAddress)
	if err != nil {
		return err
	}

	s.server = &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			path := ctx.Path()
			switch {
			case bytes.Equal(path, announcePath):
				s.announce(ctx)
			case bytes.Equal(path, scrapePath):
				s.scrape(ctx)
			default:
				ctx.NotFound()
			}
		},
		DisableKeepalive:   true,
		MaxRequestBodySize: 1024,
		GetOnly:            true,
	}

	return s.server.Serve(listener)
}

func (s *httpServer) failure(ctx *fasthttp.RequestCtx, message []byte, code int) {
	body, err := bencode.Marshal(map[string]interface{}{
		"failure reason": message,
	})
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	ctx.Error(string(body), code)
}

func New(storage storage.Storage) protocol.Server {
	return &httpServer{
		storage: storage,
	}
}

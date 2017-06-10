package http

import (
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

type event uint8

const (
	eventUnknown event = iota
	eventStarted
	eventCompleted
	eventStopped
)

type request struct {
	infoHash   []byte
	peerID     []byte
	port       uint16
	uploaded   uint64
	downloaded uint64
	left       uint64
	compact    bool
	noPeerID   bool
	event      event
	ip         []byte
	numWant    uint8
	key        []byte
	trackerID  []byte
}

func makeRequest(ctx *fasthttp.RequestCtx) (*request, error) {
	args := ctx.URI().QueryArgs()

	infoHash := args.PeekBytes(infoHashKey)
	if infoHash == nil || len(infoHash) != 20 {
		return nil, errors.New("btt: invalid info hash")
	}

	peerID := args.Peek("peer_id")
	if peerID == nil || len(peerID) != 20 {
		return nil, errors.New("btt: invalid peer id")
	}

	port, err := args.GetUint("port")
	if err != nil {
		return nil, err
	}
	if port == 0 {
		return nil, errors.New("btt: invalid port")
	}

	uploaded, err := args.GetUint("uploaded")
	if err != nil {
		return nil, err
	}

	downloaded, err := args.GetUint("downloaded")
	if err != nil {
		return nil, err
	}

	left, err := args.GetUint("left")
	if err != nil {
		return nil, err
	}

	var ip []byte
	ipRaw := args.Peek("ip")
	if len(ipRaw) != 0 {
		parts := strings.Split(string(ipRaw), ".")
		if len(parts) != 4 {
			return nil, errors.New("btt: invalid ip")
		}

		ip = make([]byte, 4)
		for i, part := range parts {
			value, err := strconv.ParseUint(part, 10, 8)
			if err != nil {
				return nil, err
			}

			ip[i] = byte(value)
		}
	} else {
		xForwardedFor := ctx.Request.Header.Peek("X-Forwarded-For")
		if ctx.RemoteIP().String() == "127.0.0.1" && len(xForwardedFor) != 0 {
			ip = net.ParseIP(string(xForwardedFor))
		} else {
			ip = ctx.RemoteIP()
		}
	}

	event := eventUnknown
	switch string(args.Peek("event")) {
	case "":
		event = eventUnknown
	case "started":
		event = eventStarted
	case "completed":
		event = eventCompleted
	case "stopped":
		event = eventStopped
	default:
		return nil, errors.New("btt: invalid event")
	}

	numWant, err := args.GetUint("numwant")
	if err != nil {
		return nil, err
	}

	if numWant == 0 || numWant > 100 {
		numWant = 50
	}

	return &request{
		infoHash:   infoHash,
		peerID:     peerID,
		port:       uint16(port),
		uploaded:   uint64(uploaded),
		downloaded: uint64(downloaded),
		left:       uint64(left),
		compact:    args.GetBool("compact"),
		noPeerID:   args.GetBool("no_peer_id"),
		event:      event,
		ip:         ip,
		numWant:    uint8(numWant),
		key:        args.Peek("key"),
	}, nil
}

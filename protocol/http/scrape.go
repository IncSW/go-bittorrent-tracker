package http

import (
	bencode "github.com/IncSW/go-bencode"
	"github.com/valyala/fasthttp"
)

func (s *httpServer) scrape(ctx *fasthttp.RequestCtx) {
	hashes := ctx.URI().QueryArgs().PeekMultiBytes(infoHashKey)
	if hashes == nil || len(hashes) == 0 {
		s.failure(ctx, failureInfoHashNotFound, fasthttp.StatusBadRequest)
		return
	}

	for _, hash := range hashes {
		if len(hash) != 20 {
			s.failure(ctx, failureInvalidInfoHash, fasthttp.StatusBadRequest)
			return
		}
	}

	files := map[string]interface{}{}
	for _, hash := range hashes {
		info, err := s.storage.GetInfo(hash)
		if err != nil {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
			return
		}

		if info == nil {
			continue
		}

		files[string(hash)] = map[string]interface{}{
			"incomplete": info.Incomplete,
			"complete":   info.Complete,
			"downloaded": info.Downloaded,
		}
	}

	if len(files) == 0 {
		s.failure(ctx, failureTorrentsNotFound, fasthttp.StatusNotFound)
		return
	}

	body, err := bencode.Marshal(map[string]interface{}{
		"files": files,
	})
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("text/plain; charset=utf8")
	ctx.Write(body)
}

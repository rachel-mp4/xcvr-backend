package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"xcvr-backend/internal/types"
)

func (h *Handler) getChannels(w http.ResponseWriter, r *http.Request) {
	limitstr := r.URL.Query().Get("limit")
	limit := 50
	if limitstr != "" {
		l, err := strconv.Atoi(limitstr)
		if err == nil {
			limit = max(min(l, 100), 1)
		}
	}
	cvs, err := h.db.GetChannelViews(limit, r.Context())
	if err != nil {
		h.serverError(w, err)
		h.logger.Printf("db.GetChannels failed! %s", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(cvs)
}

func (h *Handler) getMessages(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) resolveChannel(w http.ResponseWriter, r *http.Request) {
	handle := r.URL.Query().Get("handle")
	did := r.URL.Query().Get("did")
	rkey := r.URL.Query().Get("rkey")
	if did == "" {
		if handle == "" {
			h.badRequest(w, errors.New("did not provide did or handle"))
			return
		}
		var err error
		did, err = h.db.ResolveHandle(handle, r.Context())
		if err != nil {
			h.serverError(w, err)
			return
		}
	}
	url := fmt.Sprintf("/lrc/%s/%s/ws", did, rkey)
	rchanres := types.ResolveChannelResponse{URL: url}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(rchanres)
}

func (h *Handler) getProfileView(w http.ResponseWriter, r *http.Request) {
	handle := r.URL.Query().Get("handle")
	did := r.URL.Query().Get("did")
	if did == "" {
		if handle == "" {
			h.badRequest(w, errors.New("did not provide did or handle"))
			return
		}
		var err error
		did, err = h.db.ResolveHandle(handle, r.Context())
		if err != nil {
			h.serverError(w, err)
			return
		}
	}
	h.serveProfileView(did, handle, w, r)
}

func (h *Handler) serveProfileView(did string, handle string, w http.ResponseWriter, r *http.Request) {
	profile, err := h.db.GetProfileView(did, r.Context())
	if err != nil {
		h.notFound(w, errors.New(fmt.Sprintf("couldn't find profile for handle %s / did %s: %s", handle, did, err.Error())))
		return
	}
	profile.Handle = handle
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(profile)
}

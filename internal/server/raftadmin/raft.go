package raftadmin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/raft"
)

func Register(router *mux.Router, r *raft.Raft) {
	handler := raftHandler{
		r: r,
	}

	router.Path("/admin/raft/add_member").Methods(http.MethodPost).HandlerFunc(handler.AddMember)
}

type raftHandler struct {
	r *raft.Raft
}

func (h *raftHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name        string `json:"id"`
		Address     string `json:"address"`
		TimeoutSecs int    `json:"timeout"`
		Voter       bool   `json:"voter"`
	}{
		Voter: true,
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, fmt.Sprintf("invalid add_member request: %q", err), http.StatusBadRequest)
		return
	}

	var f raft.IndexFuture
	if data.Voter {
		f = h.r.AddVoter(raft.ServerID(data.Name), raft.ServerAddress(data.Address), 0, time.Duration(data.TimeoutSecs)*time.Second)
	} else {
		f = h.r.AddNonvoter(raft.ServerID(data.Name), raft.ServerAddress(data.Address), 0, time.Duration(data.TimeoutSecs)*time.Second)
	}

	if err := f.Error(); err != nil {
		http.Error(w, fmt.Sprintf("failed to add voter: %q", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Drofff/revsynth-server-http/service"
)

func HandleRequest(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("ERROR: unable to read request: %e", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	in := &service.SynthesiseInput{}
	err = json.Unmarshal(data, in)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("invalid input format"))
		if err != nil {
			log.Println(err)
			return
		}
	}

	out, err := service.Synthesise(in)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Println(err)
			return
		}
	}

	outJSON, err := json.Marshal(out)
	if err != nil {
		log.Printf("ERROR: unable to marshal the result: %e", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(outJSON)
	if err != nil {
		log.Println(err)
	}
}

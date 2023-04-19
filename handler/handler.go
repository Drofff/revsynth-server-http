package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Drofff/revsynth-server-http/service"
)

const (
	headerAllowOrigin  = "Access-Control-Allow-Origin"
	headerAllowMethods = "Access-Control-Allow-Methods"
	headerAllowHeaders = "Access-Control-Allow-Headers"
	headerContentType  = "Content-Type"

	allowAll = "*"

	contentTypeJSON = "application/json"
)

func addCORSHeadersPOST(w http.ResponseWriter) {
	w.Header().Add(headerAllowOrigin, allowAll)
	w.Header().Add(headerAllowMethods, http.MethodPost)
	w.Header().Add(headerAllowHeaders, allowAll)
}

func HandleRequest(w http.ResponseWriter, req *http.Request) {
	addCORSHeadersPOST(w)

	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

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
		log.Printf("WARN: marshal error: %v", err.Error())

		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("{\"msg\":\"invalid input format\"}"))
		if err != nil {
			log.Println(err)
		}
		return
	}

	out, err := service.Synthesise(in)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Println(err)
		}
		return
	}

	outJSON, err := json.Marshal(out)
	if err != nil {
		log.Printf("ERROR: unable to marshal the result: %e", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(headerContentType, contentTypeJSON)
	_, err = w.Write(outJSON)
	if err != nil {
		log.Println(err)
	}
}

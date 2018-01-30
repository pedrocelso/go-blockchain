package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/pedrocelso/go-blockchain/services/block"
)

type Message struct {
	BPM int
}

var BlockChain []block.Block

func generateBlock(oldBlock block.Block, BPM int) (*block.Block, error) {
	var err error
	t := time.Now()

	newBlock := &block.Block{}

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash

	newBlock.Hash, err = newBlock.GenerateHash()

	if err != nil {
		return nil, err
	}

	return newBlock, nil
}

func replaceChain(newBlocks []block.Block) {
	if len(newBlocks) > len(BlockChain) {
		BlockChain = newBlocks
	}
}

func run() error {
	mux := makeMuxRouter()
	s := &http.Server{
		Addr:           ":9898",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s.ListenAndServe()
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}

	w.WriteHeader(code)
	w.Write(response)
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(BlockChain, "", "  ")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}

	defer r.Body.Close()
	newBlock, err := generateBlock(BlockChain[len(BlockChain)-1], m.BPM)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	if newBlock.IsBlockValid(BlockChain[len(BlockChain)-1]) {
		newBlockchain := append(BlockChain, *newBlock)
		replaceChain(newBlockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func main() {
	go func() {
		t := time.Now()
		firstBlock := block.Block{0, t.String(), 0, "", ""}
		spew.Dump(firstBlock)
		BlockChain = append(BlockChain, firstBlock)
	}()

	log.Fatal(run())
}

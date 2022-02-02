package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	htgotts "github.com/hegedustibor/htgo-tts"
)

const audioFolder = "audio"
const fileName = "speech"
const defaultPort = 1337
const volume = 2

func main() {
	http.HandleFunc("/", serveSpeech)
	port := defaultPort
	if os.Getenv("TTSAAS_PORT") != "" {
		var err error
		port, err = strconv.Atoi(os.Getenv("TTSAAS_PORT"))
		if err != nil {
			log.Fatalf("Port environment variable set, but could not convert it to int %s\n", err)
		}
	}
	log.Printf("Starting text to speech as a service on port %d\n", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		panic(err)
	}
}

func serveSpeech(w http.ResponseWriter, r *http.Request) {
	addCORSHeader(w)
	b, _ := io.ReadAll(r.Body)
	urlParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(urlParts) < 1 {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Println("Bad request")
		return
	}
	sentence := string(b) 

	speech := htgotts.Speech{Folder: audioFolder, Language: "no"}
	speechFile, _ := speech.CreateSpeechFile(sentence,fileName)

	fileURI := speechFile
	loudFileURI := audioFolder + "/" + fileName + "LOUD" + ".mp3"
	cmdString := fmt.Sprintf("ffmpeg -y -i %s -filter:a \"volume=%d\" %s", fileURI, volume, loudFileURI)
	cmd := exec.Command("bash", "-c", cmdString)
	buf, err := cmd.Output()
	e := os.Remove(speechFile)
	if (err != nil || e != nil) {
		http.Error(w, "server error", http.StatusInternalServerError)
		log.Printf("Error increasing audio volume: %s: %s\n", err, string(buf))
		log.Printf(cmdString)
		return
	}

	time.Sleep(500 * time.Millisecond)
	http.ServeFile(w, r, loudFileURI)
}

func addCORSHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

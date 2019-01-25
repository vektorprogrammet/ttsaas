package main

import (
    "github.com/hegedustibor/htgo-tts"
    "github.com/kennygrant/sanitize"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
)

const audioFolder = "audio"
const defaultPort = 80

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
    if err := http.ListenAndServe(":" + strconv.Itoa(port), nil); err != nil {
        panic(err)
    }
}

func serveSpeech(w http.ResponseWriter, r *http.Request) {
    addCORSHeader(w)

    urlParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
    if len(urlParts) < 1 {
        http.Error(w, "bad request", http.StatusBadRequest)
        log.Println("Bad request")
        return
    }
    sentence := sanitize.BaseName(urlParts[0])

    // Save audio file to audio folder
    speech := htgotts.Speech{Folder: audioFolder, Language: "no"}
    err := speech.Speak(sentence)
    if err != nil {
        http.Error(w, "server error", http.StatusInternalServerError)
        log.Printf("Error converting text to speech: %s\n", err)
        return
    }

    http.ServeFile(w, r, audioFolder + "/" + sentence + ".mp3")
}

func addCORSHeader(w http.ResponseWriter)  {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers",
        "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

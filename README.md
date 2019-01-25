# TTSAAS
> Text To Speech As A Service

## Dependencies
Be sure to install the following go libraries:
```bash
go get github.com/hegedustibor/htgo-tts
go get github.com/kennygrant/sanitize
```

Also, the `htgo-tts` package requires `mplayer` to be installed, or it will throw an error.
```bash
sudo apt install mplayer
```


## Usage 
### Setup
You might need to create a folder named `audio` in the directory of the executable.

### Running 
Start the project with `go run main.go` or build and run the executable.

### Choosing a port
The service uses port 80 by default, but this can be configured by setting the `TTSAAS_PORT` environment variable.

## API usage
Navigating to `localhost:8000/hello%20world` will return an mp3 file with "Hello world" spoken out. The content type is `audio/mpeg`.

## HTML audio element example
The api can be used as a source for audio tags. For example:
```html
<audio controls>
  <source src="http://localhost/hello%20world" type="audio/mpeg">
</audio>
```

# inari

photo management system

## ui

### inari-web

```
make build-web
docker run -p 8080:80 --rm inari-web
http://localhost:8080/
```

## cli

### import media

```
make build-cli
docker run -v ~/photos/Camera:/inbox -v ~/photos_inari_mediastore:/mediastore -it --rm inari-cli ./inari import /inbox/IMG_20181229_140303.jpg
```

# Roadmap

## version 1

- import media
- list collections
- list collection media / view collection detail
- view media detail
- update media comment
- delete media

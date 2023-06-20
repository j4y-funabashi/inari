# inari

photo management system

## ui

### inari-web

```
make build-web
docker run -p 8080:80 --rm inari-web
docker run -v ~/photos_inari_mediastore:/tmp/inari -v ~/photos_inari_mediastore/thumbnails:/var/media/thumbnails -p 8080:80 --rm inari-web
http://localhost:8080/
```

## cli

### import media

```
make build-cli
docker run --env GOOGLE_API_KEY=${GOOGLE_API_KEY} -v ~/Downloads/jayr-phone-camera:/inbox -v ~/photos_inari_mediastore:/tmp/inari -it --rm inari-cli ./inari import /inbox/
docker run --env GOOGLE_API_KEY=${GOOGLE_API_KEY} -v /mnt/data/backup/jayr/phone/Camera/:/inbox -v /mnt/data/backup/jayr/inari_mediastore:/tmp/inari -it --rm inari-cli ./inari import /inbox/

docker run --env GOOGLE_API_KEY=${GOOGLE_API_KEY} -v /mnt/data/backup/jayr/phone/Camera/:/inbox -v /mnt/data/backup/jayr/inari_mediastore:/tmp/inari -it --rm inari-cli ./inari import /inbox/
```

# Roadmap

## version 1

- import media
- list collections
- view collection detail
- delete media

## future

- full text search captions
- update media caption
- download original media file
- publish public photos (micropub?)
- search for and add location to media
- mark as read / remove from inbox / remove from collection?
- view media detail

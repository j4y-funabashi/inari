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


docker run --env GOOGLE_API_KEY=${GOOGLE_API_KEY} -v /mnt/data/backup/jayr/phone/Camera/:/inbox/phone -v /mnt/data/backup/jayr/camera/:/inbox/camera -v /mnt/data/backup/jayr/j4y.co/:/inbox/j4y  -v /mnt/data/backup/jayr/inari_mediastore:/tmp/inari -it --rm inari-cli ./inari import /inbox/
```

# Roadmap

## version 1

- [x] import media
- [x] list collections
- [X] view collection detail
- [x] delete media
- [x] update media caption
- [x] import geo xml files
- [x] use imported locations to add lat/lng on import
- [ ] https://github.com/j4y-funabashi/inari/issues/10
- [ ] search for and add location to media
- [ ] search for and add venue to media
- [ ] import movie files

## future

- sort out homepage, list other collection types?
- generate static map when geocoding
- record import errors?
- add sizes to thumbnail objects?
- full text search captions
- download original media file
- publish public photos (micropub?)
- mark as read / remove from inbox / remove from collection?
- change collection view to current media + thumbnails
- gpx viewer
    - list of days
    - gpx day detail: map (leaflet?), stats etc

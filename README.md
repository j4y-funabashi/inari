# inari

photo management system

## ui

### inari-web

```
make build-web
docker run -p 8080:80 --rm inari-web
```

## cli

### import media

```
make build-cli
docker run -v ~/Downloads:/inbox -it --rm inari-cli ./inari import /inbox/20210102_090011_7dd977376d76f24a949838807375d831.JPG
```

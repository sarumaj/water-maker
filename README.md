# water-maker

An app to watermark PNG, JPEG and GIF files. 

## Usage

The watermark file can be overwritten with the **WATERMARK_FILE** env variable.

## Build

Reruires **gcc** installed and **GCO_ENABLED** set to `1`.

```
go build -ldflags="-s -w" . -o water-maker
```

## Screenshots
![select](doc/selection.png)

![progress](doc/progress.png)

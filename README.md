# MP4 Parser & Viewer

Experimental [MP4](https://en.wikipedia.org/wiki/MP4_file_format) parser/viewer for learning purpose.

## Intro

- It is originated from Apple [QuickTime File Format](https://en.wikipedia.org/wiki/QuickTime_File_Format) (.mov);
- then it extended to MP4;
- then it generalized into ISOBMFF, which is the basis of 3GP, JPEG 2000 as well.

```text
QuickTime File Format
    --> ISOBMFF
        --> MP4/3GP
```

### Atom/Box

https://developer.apple.com/documentation/quicktime-file-format/atoms

![atom](https://docs-assets.developer.apple.com/published/2579cfdeca4506b895b19138b92c1ae8/sample-atom~dark@2x.png)

### Basic structure

```text
some.mp4
├───ftyp -------------------> FileType
├───mdat -------------------> Media Data
├───moov -------------------> Movie
│   ├───mvhd ---------------> Movie Header
│   ├───trak ---------------> Track/Stream
│   │   ├─── tkhd ----------> Track Header
│   │   └─── mdia ----------> Media Info
│   │        └─── ...
│   └───trak
│   │   ├─── tkhd ---------->
│   │   └─── mdia ---------->
│   │        └─── ...
└───udta -------------------> Userdata Box
```

## Refs

- https://dev.to/alfg/a-quick-dive-into-mp4-57fo
- https://alfg.dev/mp4-inspector/
- https://github.com/mshafiee/mp4

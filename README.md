# Ark

Manages an archive of media files, identifying and skipping duplicates on import; it archives files by their creation date. Use at your own risk.

**Note:** it can only guarantee atomic file moves on UNIX filesystems.

## Raison d'être

Over the years I have accumulated hundreds of GBs of photos and videos, stored on a multitude of removable drives. I eventually bought a NAS for home usage and moved all files there, but there's a lot of duplication and mess because of backups taken over the years.

There are of course plenty of applications on the market that can manage an archive of media, but this looked like (and was, in fact) an interesting pet project.

## Decision log

- I'm going to consider a file to be the duplicate of another one if and only if hashing them yields the same result: no other attributes (e.g. file name, creation date, ...) are taken into account
- I'm going to compute file hashes using the [go porting](https://github.com/lukechampine/blake3) of the BLAKE3 cryptographic hash function, because of its performance
- Server and client will communicate over gRPC: this enables them to run on different machines and to leverage HTTP/2 to stream files over the network
- Instead of vanilla gRPC, I'm going to use the [connect-go](https://github.com/connectrpc/connect-go) library, primarily to experiment with it
- Clients will authenticate their requests using JWT tokens, leveraging connect-go's [interceptors](https://connect.build/docs/go/interceptors)
- I'm going to report some key metrics to [Prometheus](https://prometheus.io/docs/introduction/overview/) to measure the system's performance
- I'm going to store file metadata (e.g. path, hash, ...) in a [Redis](https://redis.io/) database

## Note: Creation date

The creation date is extracted, whenever possible, from the file's [EXIF](https://exiftool.org/TagNames/EXIF.html) header. When that is not possible (either because the file type is not supported or there is no EXIF data), the file modification time is used as a fallback.

EXIF can currently be parsed from:

- JPEG, thanks to [go-jpeg-image-structure](https://github.com/dsoprea/go-jpeg-image-structure)
- HEIC, thanks to [go-heic-exif-extractor](https://github.com/dsoprea/go-heic-exif-extractor)
- TIFF-like headers such as TIFF, CR2, and ORF using my own [tiff-parser](https://github.com/fedragon/tiff-parser)

## Components

### Server

Runs on a dedicated machine (a NAS or wherever you'd like to store your media).
Receives `UploadFile` gRPC requests from clients, archiving files by creation date. It identifies files by their pre-computed hash and skips any duplicates that may be submitted for upload.

### Client

May run on any machine having network access to the server.
Recursively walks through a directory containing media files, computing the hash of each of them and issuing `UploadFile` requests to the server. It initially only sends the file metadata: the actual file content is only sent (in chunks) if the server confirms that it's not a duplicate.

## How it works

The diagram below describes how a Client uploads files to the Server. For brevity's sake, the diagram only shows how a single file is uploaded and errors are not displayed. Any error will break the circuit.

```
                         +---------+                +---------+                +-----+ +-----+
                         | Client  |                | Server  |                | DB  | | HDD |
                         +---------+                +---------+                +-----+ +-----+
                              |                          |                        |       |
                              | Compute file hash        |                        |       |
                              |------------------        |                        |       |
                              |                 |        |                        |       |
                              |<-----------------        |                        |       |
                              |                          |                        |       |
                              | Send file metadata       |                        |       |
                              |------------------------->|                        |       |
                              |                          |                        |       |
                              |                          | Does it exist          |       |
                              |                          |----------------------->|       |
         -------------------\ |                          |                        |       |
         | alt: file exists |-|                          |                        |       |
         |------------------| |                          |                        |       |
                              |                          |                        |       |
                              |                          |                    Yes |       |
                              |                          |<-----------------------|       |
                              |                          |                        |       |
                              |      File already exists |                        |       |
                              |<-------------------------|                        |       |
                              |                          |                        |       |
                              | Skip file                |                        |       |
                              |----------                |                        |       |
                              |         |                |                        |       |
                              |<---------                |                        |       |
----------------------------\ |                          |                        |       |
| else: file does not exist |-|                          |                        |       |
|---------------------------| |                          |                        |       |
----------------------------\ |                          |                        |       |
| loop: for each file chunk |-|                          |                        |       |
|---------------------------| |                          |                        |       |
                              |                          |                        |       |
                              | Send file chunk          |                        |       |
                              |------------------------->|                        |       |
                              |                          |                        |       |
                              |                          | Store file chunk       |       |
                              |                          |-----------------       |       |
                              |                          |                |       |       |
                              |                          |<----------------       |       |
                              |                          |                        |       |
                              |                       OK |                        |       |
                              |<-------------------------|                        |       |
                 -----------\ |                          |                        |       |
                 | end loop |-|                          |                        |       |
                 |----------| |                          |                        |       |
                              |                          |                        |       |
                              |                          | Atomically write file  |       |
                              |                          |------------------------------->|
                              |                          |                        |       |
                              |                          |                        |    OK |
                              |                          |<-------------------------------|
                              |                          |                        |       |
                              |                          | Store file metadata    |       |
                              |                          |----------------------->|       |
                              |                          |                        |       |
                              |                          |                     OK |       |
                              |                          |<-----------------------|       |
                              |                          |                        |       |
                              |                       OK |                        |       |
                              |<-------------------------|                        |       |
                      ------\ |                          |                        |       |
                      | end |-|                          |                        |       |
                      |-----| |                          |                        |       |
                              |                          |                        |       |
```

### Credits

The diagram has been generated by https://weidagang.github.io/text-diagram/ using the following script:

```
object Client Server DB HDD
Client->Client: Compute file hash
Client->Server: Send file metadata
Server->DB: Does it exist
note left of Client: alt: file exists
DB->Server: Yes
Server->Client: File already exists
Client->Client: Skip file
note left of Client: else: file does not exist
note left of Client: loop: for each file chunk
Client->Server: Send file chunk
Server->Server: Store file chunk
Server->Client: OK
note left of Client: end loop
Server->HDD: Atomically write file
HDD->Server: OK
Server->DB: Store file metadata
DB->Server: OK
Server->Client: OK
note left of Client: end
```

## EXIF parsing resources

- https://exiftool.org/TagNames/EXIF.html
- http://lclevy.free.fr/cr2/
- https://github.com/lclevy/libcraw2/blob/master/docs/cr2_poster.pdf
- https://github.com/ImranAtBhimsoft/metadata-extractor

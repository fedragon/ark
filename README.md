# Ark

Manages an archive of media files, identifying and skipping duplicates on import. It archives files by their creation date.

## Components

### Server

Runs on the NAS itself.
Receives `UploadFile` requests from clients, archiving files by creation date. It identifies files by their pre-computed hash and skips any duplicates that may be submitted for upload.

### Client

May run on any machine having network access to the NAS.
Recursively walks through a directory containing media files, computing the hash of each file conforming to the configured file types and issuing `UploadFile` requests to the server.

## Upload logic

The diagram below describes how the Client uploads files to the Server. For brevity's sake, the diagram only shows how a single file is uploaded and errors are not displayed. Any error will short-circuit the whole flow.

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
                              | Upload file              |                        |       |
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
Client->Server: Upload file
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

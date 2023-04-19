# Ark

Manages an archive of media files, identifying and skipping duplicates on import. It archives files by their creation date.

## Components

### Server

Runs on the NAS itself.
Receives `UploadFile` requests from clients, archiving files by creation date. It identifies files by their pre-computed hash and skips any duplicates that may be submitted for upload.

### Client

May run on any machine having network access to the NAS.
Recursively walks through a directory containing media files, computing the hash of each file conforming to the configured file types and issuing `UploadFile` requests to the server.

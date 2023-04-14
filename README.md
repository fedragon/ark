# Ark

Manages an archive of media files, identifying and skipping duplicates on import. It archives files by their creation date.

## Components

### Server

Runs on the NAS itself and is responsible for:

- Identifying duplicates based on file hashes
- Archiving files by creation date
- Generating a static site for browsing the archive

### Client

Runs on a client machine and is responsible for:

- Sending a request to import files from a source directory, providing a list of hashes of files to be imported
- Receiving a list of non-duplicate files to be imported from the server
- Sending requests to import non-duplicate files

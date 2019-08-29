
## Installation

First make sure Go is installed. If version 1.12 is not available for your
operating system, the latest version can be gotten from:

    http://golang.org/dl

Do not use a newer version than 1.10.

Then do,

    go get github.com/paracrawl/giawarc/...

which will build and place some programs in `${HOME}/go/bin`. They
should have no runtime dependencies that stop it being moved to any machine,
it can run just fine out of an NFS directory.

## Usage

### giawarc

This program takes WARC files, cleans and preprocesses them. It can write
output compatible with [Bitextor][1]'s warc2preprocess program. It can
also write gzipped output of the form,

    Content-Location: http://edutweetoz.org/
    Content-Type: text/html
    Content-Language: en
    Content-Length: 44676
    Date: 2014-12-17T06:08:19Z
    X-WARC-Record-ID: <urn:uuid:8b505367-9444-4af7-952d-b2e6d3430bc3>
    X-WARC-Filename: SURV-20141217060818673-00010-1372~crawl421.us.archive.org~9443.warc.gz
    
    Edutweetoz | Celebrating Australian Educators Edutweetoz
    Celebrating Australian Educators
    Menu Skip to content
    About
    Catch up on what’s been happening via our storify
    ...
    
These are some headers, just like with HTTP, and then a blank line followed
by the cleaned content. Note that the Content-Type is the *original* one and
what follows here is simply text. The text is minimally processed, transformed
to UTF-8 and normalised, with excess spacing removed and newlines added where
HTML ought to put them.

Each record in the gzip file corresponds to a document in the original WARC 
file. Though it can simply be read with `gzip -dc`, the correct way is to 
read the separate gzip streamsr: there is one stream per document.

There is also an output mode that creates files like this split by language.

[1]: https://github.com/bitextor/bitextor

### giarec

This program reads files created by `giawarc`. It can extract header fields
and text bodies, optionally encoding them with Base64. In this way, a single
output file created in the `gzip` format described above can be used to
generate a dataset that is readable by `bitextor`.

### htmlstrip

This program filters and cleans HTML documents provided on the standard input
and writes them to the standard output. It uses the same mechanism as `giawarc`
and can be used to check that the output is correct.



## Installation

First make sure Go is installed. On ancient versions of Ubuntu Linux, the
recipe is:

    apt-get install software-properties-common python-software-properties
    add-apt-repository ppa:longsleep/golang-backports
    apt-get update
    apt-get install golang-go

Then do,

    go get github.com/wwaites/giawarc

which will result in the `giawarc` binary in `${HOME}/go/bin`. The binary
should have no runtime dependencies that stop it being moved to any machine,
it can run just fine out of an NFS directory.

## Usage

For now, it processes a WARC file from the command line and makes output in
the current directory. It doesn't do flags or anything like that yet. Just
testing...


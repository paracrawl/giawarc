
## Installation

First make sure Go is installed. On recent operating systems, simply
installing from the system package manager should do. On ancient versions
of Ubuntu Linux, the recipe is:

    apt-get install software-properties-common python-software-properties
    add-apt-repository ppa:gophers/archive
    apt-get update
    apt-get install golang-1.10-go
    ln -s /usr/lib/go-1.10/bin/go /usr/local/bin

Do not use a newer version than 1.10.

Then do,

    go get github.com/wwaites/giawarc

which will result in the `giawarc` binary in `${HOME}/go/bin`. The binary
should have no runtime dependencies that stop it being moved to any machine,
it can run just fine out of an NFS directory.

## Usage

For now, it processes a WARC file from the command line and makes output in
the current directory. It doesn't do flags or anything like that yet. Just
testing...


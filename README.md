# dcrseedgen

**dcrseedgen** is a desktop application that generates mnemonic
seed and converts them into hex seeds for use in decred software.

## Installation 

### Option 1: Get the binary 
**dcrseedgen** has not been released yet. For now, you have to build from source 

### Option 2: Build from source 

#### Step 1. Install Go 
* Minimum supported version is 1.11.4. Installation instructions can be found [here](https://golang.org/doc/install).
* Set `$GOPATH` environment variable and add `$GOPATH/bin` to your PATH environment variable as part of the go installation process.

#### Step 2. Clone this repo 
It is conventional to clone to $GOPATH, but not necessary.
```bash
git clone https://github.com/raedahgroup/dcrseedgen $GOPATH/src/github.com/raedahgroup/dcrseedgen
```

#### Step 3. Build the source code
* If you cloned to $GOPATH, set the `GO111MODULE=on` environment variable before building.
Run `export GO111MODULE=on` in terminal (for Mac/Linux) or `setx GO111MODULE on` in command prompt for Windows.
* `cd` to the cloned project directory and run `go build` or `go install`.
Building will place the `dcrseedgen` binary in your working directory while install will place the binary in $GOPATH/bin.

## Contributing 

See the CONTRIBUTING.md file for details. Here's an overview:

1. Fork this repo to your github account
2. Before starting any work, ensure the master branch of your forked repo is even with this repo's master branch
2. Create a branch for your work (`git checkout -b my-work master`)
3. Write your codes
4. Commit and push to the newly created branch on your forked repo
5. Create a [pull request](https://github.com/raedahgroup/dcrseedgen/pulls) from your new branch to this repo's master branch
- [Package Management](#package-management)
  - [GOPATH](#gopath)
  - [Go Vendor](#go-vendor)
  - [Go Mod](#go-mod)

# Package Management

## GOPATH

bin: contains compiled executable binary files
pkg: contains compiled .a files
src: contains project source code, which can be written by oneself or `go get`

The way to manage your own packages or other people's packages under `$GOPATH/src` is called the GOPATH mode. Under this mode, executable files generated using go install are stored in `$GOPATH/bin`

Problems:

- Unable to use specified versions of packages within a project
- Other people running your program cannot ensure that the version of the package they download is the version you expect, and using different versions of the package may cause the program to run abnormally
- Only one version of a package can be retained locally, meaning that all local projects can only use the same version of the package

## Go Vendor

> In each project, create a vendor directory, and all dependencies required by the project will only be downloaded to their own vendor directory. The dependencies between projects do not affect each other. In version 1.5 of Golang, when setting `GO15VENDOREXPERIMENT=1`, the priority of searching for vendor directory dependencies (compared to GOPATH) is improved.


The search priority for packages from high to low:

- vendor directory within the current package
- Search upwards until you find the vendor directory under src
- Search in the GOROOT directory
- Search for dependencies in GOPATH

Problems:

- If multiple projects use the same version of the same package, the package will exist in different directories on the machine, which is a waste of disk space and cannot be centrally managed for third-party packages
- If you want to open source a project, you need to include all dependencies in the project, which will increase the size of the project and make it difficult to distribute
- It is difficult to manage multiple versions of a package in the vendor directory

## Go Mod

Since v1.11, go env has added an environment variable: `GO111MODULE`, which can be used to enable or disable go mod mode.

It has three options:

- off: disables module support, and dependencies will be searched for in `GOPATH` & `vendor` when compiling
- on: enables module support, and `GOPATH` & `vendor` will be ignored during compilation, and dependencies will only be downloaded based on go.mod
- auto: automatically enables module support when the project is outside of `$GOPATH/src` and the root directory of the project has a go.mod file
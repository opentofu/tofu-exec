# This repository is a work in progress

This repository is currently not usable and requires a large set of changes to make it so.
We are going to update this repository solely to use it in [`tofu-ls`](https://github.com/opentofu/tofu-ls).
This document will be updated as part of those changes and once deemed ready this notice will be removed. We still do not recommend using this library outside of [`tofu-ls`](https://github.com/opentofu/tofu-ls). [Link to the Issue](https://github.com/opentofu/opentofu/issues/2455#issuecomment-2858320418)


[![PkgGoDev](https://pkg.go.dev/badge/github.com/opentofu/tofu-exec)](https://pkg.go.dev/github.com/opentofu/tofu-exec)

# tofu-exec

A Go module for constructing and running [OpenTofu](https://opentofu.org/) CLI commands. Structured return values use the data types defined in [terraform-json](https://github.com/hashicorp/terraform-json).

Currently, this library is built and maintained for a few specific uses in other OpenTofu projects (such as [`tofu-ls`](https://github.com/opentofu/tofu-ls), and is not intended for general purpose use.

## Go compatibility

This library is built in Go, and uses the [support policy](https://golang.org/doc/devel/release.html#policy) of Go as its support policy. At least, the two latest major releases of Go are supported by tofu-exec.

Currently, that means Go **1.18** or later must be used.

## Usage

The `Tofu` struct must be initialised with `NewTofu(workingDir, execPath)`.

Top-level `tofu` commands each have their own function, which will return either `error` or `(T, error)`, where `T` is a `terraform-json` type.


### Example

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/opentofu/tofudl"
	"github.com/opentofu/tofu-exec/tfexec"
)

// Temporary install and execution for tofu using tofudl and tofu-exec
func main() {
	// Creating temporary directory to put our binary in
	tempDir, err := os.MkdirTemp("", "tofuinstall")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	dl, err := tofudl.New()
	if err != nil {
		log.Fatalf("error when instantiating tofudl %s", err)
	}

	// Downloading and writing tofu binary v1.9.1
	ver := tofudl.Version("1.9.1")
	opts := tofudl.DownloadOptVersion(ver)
	binary, err := dl.Download(context.TODO(), opts)
	if err != nil {
		log.Fatalf("error when downloading %s", err)
	}

	execPath := filepath.Join(tempDir, "tofu")
	// Windows executable case
	if runtime.GOOS == "windows" {
		execPath += ".exe"
	}
	if err := os.WriteFile(execPath, binary, 0755); err != nil {
		log.Fatalf("error when writing the file %s: %s", execPath, err)
	}

	// workingDir := "/path/to/working/dir"
	tf, err := tfexec.NewTofu(workingDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		log.Fatalf("error running Show: %s", err)
	}

	fmt.Println(state.FormatVersion) // "1.0"
}
```

## Testing Tofu binaries

The tofu-exec test suite contains end-to-end tests which run realistic workflows against a real Tofu binary using `tfexec.Tofu{}`.

To run these tests with a local Tofu binary, set the environment variable `TFEXEC_E2ETEST_TERRAFORM_PATH` to its path and run:
```sh
go test -timeout=20m ./tfexec/internal/e2etest
```

For more information on tofu-exec's test suite, please see Contributing below.

## Contributing

Please see [CONTRIBUTING.md](./CONTRIBUTING.md).

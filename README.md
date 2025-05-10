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

Top-level OpenTofu commands each have their own function, which will return either `error` or `(T, error)`, where `T` is a `terraform-json` type.


### Example

```go
// TODO: update this example once we have a final version of the API with `tofudl` library setup
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/opentofu/tofu-exec/tfexec"
)

func main() {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.0.6")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	workingDir := "/path/to/working/dir"
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

	fmt.Println(state.FormatVersion) // "0.1"
}
```

## Testing OpenTofu binaries

The tofu-exec test suite contains end-to-end tests which run realistic workflows against a real OpenTofu binary using `tfexec.Tofu{}`.

To run these tests with a local OpenTofu binary, set the environment variable `TFEXEC_E2ETEST_TERRAFORM_PATH` to its path and run:
```sh
go test -timeout=20m ./tfexec/internal/e2etest
```

For more information on tofu-exec's test suite, please see Contributing below.

## Contributing

Please see [CONTRIBUTING.md](./CONTRIBUTING.md).

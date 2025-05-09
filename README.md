# This repository is a work in progress

This repository is currently not usable and requires a large set of changes to make it so.
We are going to update this repository solely to use it in [`tofu-ls`](https://github.com/opentofu/tofu-ls).
This document will be updated as part of those changes. We do not recommend using this library outside of [`tofu-ls`](https://github.com/opentofu/tofu-ls). [Link to the Issue](https://github.com/opentofu/opentofu/issues/2455#issuecomment-2858320418)


[![PkgGoDev](https://pkg.go.dev/badge/github.com/opentofu/tofu-exec)](https://pkg.go.dev/github.com/opentofu/tofu-exec)

# tofu-exec

A Go module for constructing and running [Terraform](https://terraform.io) CLI commands. Structured return values use the data types defined in [terraform-json](https://github.com/hashicorp/terraform-json).

The [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) is the canonical Go interface for Terraform plugins using the gRPC protocol. This library is intended for use in Go programs that make use of Terraform's other interface, the CLI. Importing this library is preferable to importing `github.com/hashicorp/terraform/command`, because the latter is not intended for use outside Terraform Core.

While tofu-exec is already widely used, please note that this module is **not yet at v1.0.0**, and that therefore breaking changes may occur in minor releases.

We strictly follow [semantic versioning](https://semver.org).

## Go compatibility

This library is built in Go, and uses the [support policy](https://golang.org/doc/devel/release.html#policy) of Go as its support policy. The two latest major releases of Go are supported by tofu-exec.

Currently, that means Go **1.18** or later must be used.

## Usage

The `Terraform` struct must be initialised with `NewTerraform(workingDir, execPath)`.

Top-level Terraform commands each have their own function, which will return either `error` or `(T, error)`, where `T` is a `terraform-json` type.


### Example


```go
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
	tf, err := tfexec.NewTerraform(workingDir, execPath)
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

## Testing Terraform binaries

The tofu-exec test suite contains end-to-end tests which run realistic workflows against a real Terraform binary using `tfexec.Terraform{}`.

To run these tests with a local Terraform binary, set the environment variable `TFEXEC_E2ETEST_TERRAFORM_PATH` to its path and run:
```sh
go test -timeout=20m ./tfexec/internal/e2etest
```

For more information on tofu-exec's test suite, please see Contributing below.

## Contributing

Please see [CONTRIBUTING.md](./CONTRIBUTING.md).

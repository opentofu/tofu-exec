[![PkgGoDev](https://pkg.go.dev/badge/github.com/hashicorp/tofu-exec)](https://pkg.go.dev/github.com/hashicorp/tofu-exec)

# tofu-exec

> [!WARNING]
> tofu-exec is maintenance only! The OpenTofu team only reviews PRs on this repository, but does not perform active development or fix bugs in this repository. You may want to take a look at [TofuDL](https://github.com/opentofu/tofudl) instead.

A Go module for constructing and running [OpenTofu](https://opentofu.org) CLI commands. Structured return values use the data types defined in [terraform-json](https://github.com/hashicorp/terraform-json).

The [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) is the canonical Go interface for OpenTofu plugins using the gRPC protocol. This library is intended for use in Go programs that make use of OpenTofu's other interface, the CLI. Importing this library is preferable to importing `github.com/opentofu/opentofu/command`, because the latter is not intended for use outside OpenTofu Core.

> [!NOTE]
> tofu-exec supports OpenTofu versions 1.7.0 and up. Version 1.6.0 and previous Terraform versions may work, but are unsupported.

## Usage

The `OpenTofu` struct must be initialised with `NewOpenTofu(workingDir, execPath)`. You will need to provide it with an OpenTofu binary, which you can obtain using [TofuDL](https://github.com/opentofu/tofudl). 

Top-level OpenTofu commands each have their own function, which will return either `error` or `(T, error)`, where `T` is a `terraform-json` type.

### Example

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/hashicorp/tofu-exec/tofuexec"
	"github.com/opentofu/tofudl"
)

func main() {
	downloader, err := tofudl.New()
	if err != nil {
		log.Fatalf("failed to initialize TofuDL: %v", err)
	}
	binaryContents, err := downloader.Download(context.Background())
	if err != nil {
		log.Fatalf("failed to download OpenTofu: %v", err)
	}
	binName := "tofu"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	workingDir := "/path/to/working/dir"
	binPath := path.Join(workingDir, binName)
	if err := os.WriteFile(binPath, binaryContents, 0777); err != nil {
		log.Fatalf("failed to write %s: %v", binPath, err)
	}

	tofu, err := tofuexec.NewOpenTofu(workingDir, binPath)
	if err != nil {
		log.Fatalf("error running NewOpenTofu: %s", err)
	}

	err = tofu.Init(context.Background(), tofuexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	state, err := tofu.Show(context.Background())
	if err != nil {
		log.Fatalf("error running Show: %s", err)
	}

	fmt.Println(state.FormatVersion) // "0.1"
}
```

## Testing OpenTofu binaries

The tofu-exec test suite contains end-to-end tests which run realistic workflows against a real OpenTofu binary using `tofuexec.OpenTofu{}`.

To run these tests with a local OpenTofu binary, set the environment variable `TOFUEXEC_E2ETEST_TOFU_PATH` to its path and run:
```sh
go test -timeout=20m ./tofuexec/internal/e2etest
```

For more information on tofu-exec's test suite, please see Contributing below.

## Contributing

Please see [CONTRIBUTING.md](./CONTRIBUTING.md).

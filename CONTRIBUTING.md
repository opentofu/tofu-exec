# Contributing to tofu-exec

While tofu-exec is already widely used, please note that this module is **not yet at v1.0.0**, and that therefore breaking changes may occur in minor releases.

We strictly follow [semantic versioning](https://semver.org).

## Repository structure

Three packages comprise the public API of the tofu-exec Go module:

### `tofuexec`

Package `github.com/opentofu/tofu-exec/tofuexec` exposes functionality for constructing and running OpenTofu CLI commands. Structured return values use the data types defined in the [hashicorp/terraform-json](https://github.com/hashicorp/terraform-json) package.

#### Adding a new OpenTofu CLI command to `tofuexec`

Each OpenTofu CLI first- or second-level subcommand (e.g. `tofu refresh`, or `tofu workspace new`) is implemented in a separate Go file. This file defines a public function on the `OpenTofu` struct, which consumers use to call the CLI command, and a private `*Cmd` version of this function which returns an `*exec.Cmd`, in order to facilitate unit testing (see Testing below).

For example:
```go
func (tf *OpenTofu) Refresh(ctx context.Context, opts ...RefreshCmdOption) error {
	cmd, err := tf.refreshCmd(ctx, opts...)
	if err != nil {
		return err
	}
	return tf.runTerraformCmd(cmd)
}

func (tf *OpenTofu) refreshCmd(ctx context.Context, opts ...RefreshCmdOption) (*exec.Cmd, error) {
	...
  	return tf.buildTerraformCmd(ctx, mergeEnv, args...), nil
}
```

Command options are implemented using the functional variadic options pattern. For further reading on this pattern, please see [Functional options for friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis) by Dave Cheney.

## Testing

We aim for full test coverage of all Terraform CLI commands implemented in `tofuexec`, with as many combinations of command-line options as possible. New command implementations will not be merged without both unit and end-to-end tests.

### Environment variables

The following environment variables can be set during testing:

 - `TOFUEXEC_E2ETEST_VERSIONS`: When set to a comma-separated list of version strings, this overrides the default list of Terraform versions for end-to-end tests, and runs those tests against only the versions supplied.
 - `TOFUEXEC_E2ETEST_TOFU_PATH`: When set to the path of a valid local Terraform executable, only tests appropriate to that executable's version are run. No other versions are downloaded or run. Note that this means that tests using `runTestWithVersions()` will only run if the test version matches the local executable exactly.

If both of these environment variables are set, `TOFUEXEC_E2ETEST_TOFU_PATH` takes precedence, and any other versions specified in `TOFUEXEC_E2ETEST_VERSIONS` are ignored.

### Unit tests

Unit tests live alongside command implementations in `tofuexec/`. A unit test asserts that the *string* version of the `exec.Cmd` returned by the `*Cmd` function (e.g. `refreshCmd`) is as expected. Minimally, commands must be tested with no options passed ("defaults"), and with all options set to non-default values. The `assertCmd()` helper can be used for this purpose. Please see `tofuexec/init_test.go` for a reasonable starting point.

### End-to-end tests

End-to-end tests test `TofuDL` in conjunction with `tofuexec`, using the former to install OpenTofu binaries and exercising the latter in as many combinations as possible, after real-world use cases.

By default, each test is run against the latest patch versions of all OpenTofu minor version releases, starting at 0.11. Copy an existing test and use the `runTest()` helper for this purpose.

#### Testing behaviour that differs between OpenTofu versions

Subject to [compatibility guarantees](https://opentofu.org/docs/language/v1-compatibility-promises/), each new version of OpenTofu CLI may:
 - Add a command or flag not previously present
 - Remove a command or flag
 - Change stdout or stderr output
 - Change the format of output files, e.g. the state file
 - Change a command's exit code
 
These and any other differences between versions should be specified in test assertions.

If the command implemented differs in any way between OpenTofu versions (e.g. a flag is added or removed, or the subcommand does not exist in earlier versions), use `t.Skip()` directives and version checks to adapt test behaviour as appropriate. For example:
https://github.com/opentofu/tofu-exec/blob/d0cb3efafda90dd47bbfabdccde3cf7e45e0376d/tfexec/internal/e2etest/validate_test.go#L15-L23

The `runTestWithVersions()` helper can be used to run tests against specific OpenTofu versions. This should be used only alongside a test using `runTest()` to cover the remaining past and future versions.

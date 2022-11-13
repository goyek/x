# Changelog

All notable changes to this library are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this library adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
as well as to [Module version numbering](https://go.dev/doc/modules/version-numbers).

## [Unreleased](https://github.com/goyek/goyek/compare/v0.1.1...HEAD)

### Added

- Add `color.NoColor` function which disables colorizing the output.
- Add `-no-color` flag to `boot.Main` that disables colorizing the output.

## [0.1.1](https://github.com/goyek/goyek/releases/tag/v0.1.1) - 2022-11-06

This release bumps `goyek` to `2.0.0-rc.9`.

## [0.1.0](https://github.com/goyek/goyek/releases/tag/v0.1.0) - 2022-11-01

This release primarily adds the `boot.Main` and `cmd.Exec` functions.

### Added

- Add `boot.Main` function which is an extension of `goyek.Main` with some
  out-of-the-box configuration and `flag` support.
- Add `cmd.Exec` function that runs commands in a Shell-like way.
- Add `cmd.Dir` option that sets the working directory.
- Add `cmd.Env` option that sets an environment variable.
- Add `cmd.Stdin` option that sets the standard input.
- Add `cmd.Stdout` option that sets the standard output.
- Add `cmd.Stderr` option that sets the standard error.
- Add `color.ReportStatus` which is an extension of `middleware.ReportStatus`.
- Add `color.CodeLineLogger` which is an extension of `goyek.CodeLineLogger`.

<!-- markdownlint-configure-file
{
  "MD024": {
    "siblings_only": true
  }
}
-->

# Changelog

All notable changes to this library are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this library adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
as well as to [Module version numbering](https://go.dev/doc/modules/version-numbers).

## [0.4.0](https://github.com/goyek/x/releases/tag/v0.4.0) - 2025-11-25

### Changed

- **BREAKING**: Change `boot.Main` command line syntax from `[flags] [--] [tasks]`
  to `[tasks] [flags] [--] [args]` to align with `goyek` v3 recommendation.
- Bump `github.com/goyek/goyek` to `3.0.0`.
- Bump other dependencies.

### Remove

- Drop support for Go 1.22.

## [0.3.0](https://github.com/goyek/x/releases/tag/v0.3.0) - 2025-03-11

### Changed

- Bump dependencies.

### Remove

- Drop support for Go 1.21.

## [0.2.0](https://github.com/goyek/x/releases/tag/v0.2.0) - 2024-08-09

### Added

- Add `color.ReportFlow` which is an extension of `middleware.ReportFlow`.
- Add `otelgoyek.Middleware` and `otelgoyek.ExecutorMiddleware` which add
  OpenTelemetry tracing instrumentation.

### Changed

- Bump `github.com/goyek/goyek` to `2.2.0`.
- Bump `github.com/fatih/color` to `1.17.0`.

### Remove

- Drop support for Go 1.17, 1.18, 1.19, 1.20.

## [0.1.7](https://github.com/goyek/x/releases/tag/v0.1.7) - 2024-01-17

### Added

- `boot.Main` buffers the output from parallel tasks to not have mixed output
  from parallel tasks execution.

### Changed

- Bump `github.com/goyek/goyek` to `2.1.0`.

## [0.1.6](https://github.com/goyek/x/releases/tag/v0.1.6) - 2023-11-13

### Changed

- Bump `github.com/fatih/color` to `1.16.0`.

## [0.1.5](https://github.com/goyek/x/releases/tag/v0.1.5) - 2023-02-08

### Changed

- Bump `github.com/goyek/goyek` to `2.0.0`.
- Bump `github.com/fatih/color` to `1.14.1`.

## [0.1.4](https://github.com/goyek/x/releases/tag/v0.1.4) - 2022-11-24

### Changed

- Bump `github.com/goyek/goyek` to `2.0.0-rc.12`.

## [0.1.3](https://github.com/goyek/goyek/releases/tag/v0.1.3) - 2022-11-13

### Fixed

- Fix flag usage descriptions used in `boot.Main`.

## [0.1.2](https://github.com/goyek/x/releases/tag/v0.1.2) - 2022-11-13

### Added

- Add `color.NoColor` function which disables colorizing the output.
- Add `-no-color` flag to `boot.Main` that disables colorizing the output.

## [0.1.1](https://github.com/goyek/x/releases/tag/v0.1.1) - 2022-11-06

This release bumps `goyek` to `2.0.0-rc.9`.

## [0.1.0](https://github.com/goyek/x/releases/tag/v0.1.0) - 2022-11-01

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

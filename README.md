# goyek x packages

[![GitHub Release](https://img.shields.io/github/v/release/goyek/x)](https://github.com/goyek/x/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/goyek/x.svg)](https://pkg.go.dev/github.com/goyek/x)
[![go.mod](https://img.shields.io/github/go-mod/go-version/goyek/x)](go.mod)
[![LICENSE](https://img.shields.io/github/license/goyek/x)](LICENSE)
[![Build Status](https://img.shields.io/github/workflow/status/goyek/x/build)](https://github.com/goyek/x/actions?query=workflow%3Abuild+branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/goyek/x)](https://goreportcard.com/report/github.com/goyek/x)

⭐ `Star` this repository if you find it valuable and worth maintaining.

## Description

This repository contains supplemental packages for [`goyek`](https://github.com/goyek/goyek)
which mainly offer convenience.

Packages in this repository depend on additional libraries
and require a newer version of Go than [`goyek`](https://github.com/goyek/goyek).
See [`go.mod`](go.mod) for details.

Package [`boot`](https://pkg.go.dev/github.com/goyek/x/boot)
contains an extension of `Flow.Main` which additionally
defines flags and configures the flow in a convenient way.

Package [`cmd`](https://pkg.go.dev/github.com/goyek/x/cmd)
offers functions for running programs in a Shell-like way.

Package [`color`](https://pkg.go.dev/github.com/goyek/x/color)
contains goyek features which additionally have colors.

## Example

See [build](build) which is this repository's own build pipeline (dogfooding).

## License

**goyek/x** is licensed under the terms of the [MIT license](LICENSE).

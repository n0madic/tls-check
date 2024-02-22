# TLS Version Checker

TLS Version Checker is a Go-based CLI tool that checks server support for various TLS versions and evaluates TLS certificate validity, including expiration checks. Simple, efficient, and requires no external dependencies.

## Features

* Checks support for TLSv1.0, TLSv1.1, TLSv1.2, and TLSv1.3.
* Indicates whether a TLS version is considered deprecated.
* Verifies the server's TLS certificate for expiration.
* Provides detailed information on certificate validity.

## Installation

To install the TLS Version Checker, you must have Go installed on your system. If you don't have Go installed, you can download it from the official website.

```bash
go install github.com/n0madic/tls-check@latest
```

## Usage

Run the program from the command line, specifying one or more servers (with optional port numbers) to check:

```bash
tls-check server1.com server2.com:8443
```

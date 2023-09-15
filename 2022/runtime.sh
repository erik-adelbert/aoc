#!/bin/bash

#shellcheck disable=SC2046  # cmdline must split
go run ./runtime/main.go $(find . -name time.md)

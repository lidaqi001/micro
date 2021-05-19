#!/usr/bin/env sh

go run ./examples/services/sing/main.go &
go run ./examples/services/speak/main.go &
go run ./examples/services/listen/main.go &
go run ./examples/services/async/main.go &
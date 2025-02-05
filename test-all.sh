#!/bin/bash

go test ./controller/... -cover -coverprofile=coverage.out

go tool cover -html=coverage.out -o coverage.html
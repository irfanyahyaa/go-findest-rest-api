#!/bin/bash

go test ./controllers/... -cover -coverprofile=coverage.out

go tool cover -html=coverage.out -o coverage.html
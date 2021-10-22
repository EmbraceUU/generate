#!/bin/bash

CGO_ENABLED=0 GOOS=windows go build -ldflags '-w -s'
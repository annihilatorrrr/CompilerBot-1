#!/bin/bash
echo "Renaming your file to main.go..."
mv /build/*.go /build/main.go
echo "Compiling your file..."
go build /build/main.go
echo "Running your file..."
./main

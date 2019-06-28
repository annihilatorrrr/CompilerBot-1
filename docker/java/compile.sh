#!/bin/bash
echo "Renaming your file to main.java..."
mv /build/*.java /build/main.java
echo "Compiling your file..."
javac /build/main.java
echo "Running your file..."
java -cp /build main

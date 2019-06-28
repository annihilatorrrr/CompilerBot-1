#!/bin/bash
echo "Renaming your file to main.js..."
mv /build/*.js /build/main.js
echo "Running your file..."
node /build/main.js

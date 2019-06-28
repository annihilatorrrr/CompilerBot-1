#!/bin/bash
echo "Renaming your file to main.py..."
mv /build/*.py /build/main.py
echo "Running your file..."
python /build/main.py

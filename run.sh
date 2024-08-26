#!/bin/bash

# Default Go command
cmd="./tmp/main"

# Append any arguments passed to the script
for arg in "$@"
do
    cmd="$cmd $arg"
done

# Run the command
echo "Running: $cmd"
$cmd

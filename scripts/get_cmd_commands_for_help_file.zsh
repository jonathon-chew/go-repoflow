#! /usr/bin/env bash

# Copy to clipboard the arguments that can be found in the cmd file in order to help build out the help function
grep "case" cmd/cmd.go \
  | sed -E 's/.*case[[:space:]]"([^"]*)".*/aphrodite.PrintBold("cyan", "\1")/' 

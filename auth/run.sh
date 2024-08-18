#!/bin/bash

BIN=./server

trap "echo 'Cleaning up...'; rm $BIN; echo '$BIN removed.'; exit" INT
$BIN
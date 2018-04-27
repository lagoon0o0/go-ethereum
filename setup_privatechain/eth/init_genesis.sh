#!/usr/bin/env bash
mygeth=$1

echo "mygeth: ${mygeth}"


echo "initializing node 01..."
$mygeth --datadir ./01/ init ./genesis.json

echo "initializing node 02..."
$mygeth --datadir ./02/ init ./genesis.json

echo "initializing node 03..."
$mygeth --datadir ./03/ init ./genesis.json
#!/bin/sh
seq 1 5 | while read i; do 
  (nohup ./client &)
done
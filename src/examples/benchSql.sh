#!/bin/bash
for i in 20 40 80 160 320 640; do go run sqlite-performacetest.go 10000 $i; done | grep seconds

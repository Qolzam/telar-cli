#!/bin/sh

for f in bin/telar*; do shasum -a 256 $f > $f.sha256; done

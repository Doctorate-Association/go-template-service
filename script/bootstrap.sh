#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=userservice
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}

#!/usr/bin/env bash
resp=$(curl http://0.0.0.0:9090/ping 2>>/dev/null)
expect="pong"
if [ "$resp" = "$expect" ]; then
    echo "success"
    exit 0
else
    echo "fail"
    exit 1
fi

#!/usr/bin/env bash
ENV=dev DB=1 RELOAD=0 setsid ./admin > adming_log 2>&1 &
echo 'relaxadmin started, pid=${pid:=$!}, status=$?'
echo "$pid" > relaxadmin.pid

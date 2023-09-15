#!/bin/bash

function check_process() {
    if [ "$1" = "" ]; then
        return 1
    fi

    PROCESS_NUM=$(ps -ef | grep "$1" | grep -v "grep" | wc -l)
    if [ $PROCESS_NUM -ge 1 ]; then
        return 0
    else
        return 1
    fi
}

CURPATH=$(cd "$(dirname "$0")"; pwd)
cd $CURPATH

chmod +x ./dbip
chmod +x ./log.sh
source ./log.sh

while [ 1 ]; do
    check_process "dbip"
    Check_RET=$?
    if [ $Check_RET -eq 1 ]; then
        log_warn "service abnormal"
        nohup ./dbip -addr :29952 -mmdb ./download/ipcc.mmdb >>info.log 2>>error.log &
    else
        log_info "service normal"
    fi
    sleep 60
done

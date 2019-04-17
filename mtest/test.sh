#!/bin/sh

TARGET="$1"

fin() {
    chmod 600 ./mtest_key
    echo "-------- host1: ep-agent log"
    ./mssh ubuntu@${HOST1} sudo journalctl -u ep-agent.service --no-pager
    echo "-------- host2: ep-agent log"
    ./mssh ubuntu@${HOST2} sudo journalctl -u ep-agent.service --no-pager
    echo "-------- host3: ep-agent log"
    ./mssh ubuntu@${HOST3} sudo journalctl -u ep-agent.service --no-pager
}
trap fin INT TERM HUP 0

$GINKGO -v -focus="${TARGET}" $SUITE_PACKAGE
RET=$?

exit $RET

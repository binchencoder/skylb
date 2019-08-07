#! /bin/bash

set +e

SERVICE_BIN="skylb-dashboard"

PIDS=$(pgrep ^skylb-dashboard)
if [ $? -ne 0 ]; then
	echo "INFO: no running ${SERVICE_BIN} found!"
	exit 0
fi

echo -e "Stopping ${SERVICE_BIN} ...\c"
for PID in $PIDS; do
	kill ${PID} >/dev/null 2>&1
done

while [ true ]; do
	echo -e ".\c"

	IDS=$(pgrep ^skylb-dashboard)
	if [ $? -ne 0 ]; then
		echo
		echo "PID: $PIDS"
		echo "Done, service ${SERVICE_BIN} was stopped."
		exit 0
	fi

	sleep 1
done

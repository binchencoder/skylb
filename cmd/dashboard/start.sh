#!/bin/sh

set -e

export ETCDCTL_API=3
ETCD_ENDPOINTS=http://192.168.64.176:2381,http://192.168.64.176:2383,http://192.168.64.176:2385

SKYLB_DASHBOARD_HOST="localhost" # Ops: update it.
SKYLB_DASHBOARD_PORT="8050"

# Login with LDAP.
LDAP_ENDPOINT="ldap.eff.com:389"

set +e

DEBUG=false
if [ ${LDAP_ENDPOINT} == "" ]; then
	DEBUG=true
	echo "WARN: in debug mode, the app does not do authentication."
	echo "WARN: don't run in debug mode for production."
fi

DIR=$(
	cd $(dirname $0)
	pwd | awk -F'/bin' '{print $1}'
)

echo "app home dir: $DIR"
cd "${DIR}"

SERVICE_BIN="skylb-dashboard"

SERVICE_LOG_DIR="${DIR}/logs"
if [ ! -e ${SERVICE_LOG_DIR} ]; then
	mkdir -p "${SERVICE_LOG_DIR}"
fi

STDOUT_FILE=${SERVICE_LOG_DIR}/stdout.log

# Check if the service is running.

PIDS=$(pgrep ^skylb-dashboard$)
if [ $? -eq 0 ]; then
	echo "ERROR: The service ${SERVICE_BIN} is running!"
	echo "PID: $PIDS"
	exit 1
fi

# Start service.

cd ${DIR}/bin/
chmod +x ${DIR}/bin/*
nohup ${DIR}/bin/skylb-dashboard \
	-host-port="${SKYLB_DASHBOARD_HOST}:${SKYLB_DASHBOARD_PORT}" \
	-cert-dir="${DIR}/certs" \
	-static-dir="${DIR}/static" \
	-db-conf="${DIR}/conf/database.conf" \
	-debug-mode=${DEBUG} \
	-ldap-endpoint="${LDAP_ENDPOINT}" \
	-etcd-endpoints="${ETCD_ENDPOINTS}" \
	-log_dir="${SERVICE_LOG_DIR}" \
	-v 2 \
	-stderrthreshold INFO \
	>${STDOUT_FILE} 2>&1 &

# Verify service start ok.

echo "sleep 5 seconds to wait for service start ..."
sleep 5

PIDS=$(pgrep ^skylb-dashboard$)
if [ $? -ne 0 ]; then
	echo "ERROR: The service ${SERVICE_BIN} failed to start!"
	exit 1
fi

echo "Service started with PID: $PIDS"

# The service may being blocked at connecting to database.
# Check if port is being listening.

line=$(netstat -alnt | grep -e [\\.\:]${SKYLB_DASHBOARD_PORT} | grep -e LISTEN)
if [ $? -eq 0 ]; then
	echo "Service ${SERVICE_BIN} started."
else
	echo "ERROR: port ${SKYLB_DASHBOARD_PORT} is not bound. this means the service is not fully started right now, check log files!"
	exit 1
fi

#!/usr/bin/env bash

set -euo pipefail

APP_NAME="checkin"
DEPLOY_DIR="/opt/neworld_check-in"
SERVICE_PATH="/etc/systemd/system/${APP_NAME}.service"
TIMER_PATH="/etc/systemd/system/${APP_NAME}.timer"

require_root() {
  if [[ "${EUID}" -ne 0 ]]; then
    echo "error: please run this script as root" >&2
    exit 1
  fi
}

require_systemd() {
  if ! command -v systemctl >/dev/null 2>&1; then
    echo "error: systemctl is not available; this script requires systemd" >&2
    exit 1
  fi
}

stop_and_disable_units() {
  systemctl disable --now "${APP_NAME}.timer" >/dev/null 2>&1 || true
  systemctl stop "${APP_NAME}.service" >/dev/null 2>&1 || true
}

remove_files() {
  rm -f "${SERVICE_PATH}" "${TIMER_PATH}"
  rm -rf "${DEPLOY_DIR}"
}

reload_systemd() {
  systemctl daemon-reload
  systemctl reset-failed >/dev/null 2>&1 || true
}

print_done() {
  cat <<EOF
Uninstalled ${APP_NAME}
Removed:
  ${SERVICE_PATH}
  ${TIMER_PATH}
  ${DEPLOY_DIR}
EOF
}

main() {
  require_root
  require_systemd
  stop_and_disable_units
  remove_files
  reload_systemd
  print_done
}

main "$@"

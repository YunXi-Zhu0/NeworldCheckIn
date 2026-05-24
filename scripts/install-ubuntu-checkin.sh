#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

APP_NAME="checkin"
DEPLOY_DIR="/opt/neworld_check-in"
BINARY_SOURCE="${REPO_ROOT}/checkin"
CONFIG_SOURCE="${REPO_ROOT}/config.yaml"
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

require_source_files() {
  if [[ ! -x "${BINARY_SOURCE}" ]]; then
    echo "error: compiled binary not found: ${BINARY_SOURCE}" >&2
    echo "build it first with: go build -o checkin ./cmd/checkin" >&2
    exit 1
  fi

  if [[ ! -f "${CONFIG_SOURCE}" ]]; then
    echo "error: config file not found: ${CONFIG_SOURCE}" >&2
    echo "copy config.yaml.example to config.yaml and fill in your credentials first" >&2
    exit 1
  fi
}

install_files() {
  install -d -m 0755 "${DEPLOY_DIR}"
  install -m 0755 "${BINARY_SOURCE}" "${DEPLOY_DIR}/${APP_NAME}"
  install -m 0600 "${CONFIG_SOURCE}" "${DEPLOY_DIR}/config.yaml"
}

install_service() {
  cat >"${SERVICE_PATH}" <<EOF
[Unit]
Description=Neworld daily check-in
After=network-online.target
Wants=network-online.target

[Service]
Type=oneshot
WorkingDirectory=${DEPLOY_DIR}
ExecStart=${DEPLOY_DIR}/${APP_NAME}
EOF
}

install_timer() {
  cat >"${TIMER_PATH}" <<EOF
[Unit]
Description=Run Neworld daily check-in at a random time between 06:00 and 09:00

[Timer]
OnCalendar=*-*-* 06:00:00
RandomizedDelaySec=3h
Persistent=true
Unit=${APP_NAME}.service

[Install]
WantedBy=timers.target
EOF
}

enable_timer() {
  systemctl daemon-reload
  systemctl enable --now "${APP_NAME}.timer"
}

print_next_steps() {
  cat <<EOF
Installed ${APP_NAME} to ${DEPLOY_DIR}
Enabled timer: ${APP_NAME}.timer

Useful commands:
  systemctl status ${APP_NAME}.timer
  systemctl list-timers ${APP_NAME}.timer
  systemctl start ${APP_NAME}.service
  journalctl -u ${APP_NAME}.service -n 50 --no-pager
EOF
}

main() {
  require_root
  require_systemd
  require_source_files
  install_files
  install_service
  install_timer
  enable_timer
  print_next_steps
}

main "$@"

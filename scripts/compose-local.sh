#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR/infra"

export XDG_RUNTIME_DIR="${XDG_RUNTIME_DIR:-/run/user/$(id -u)}"
export DOCKER_HOST="unix://${XDG_RUNTIME_DIR}/podman/podman.sock"

mkdir -p "${XDG_RUNTIME_DIR}/podman"
if [[ ! -S "${XDG_RUNTIME_DIR}/podman/podman.sock" ]]; then
  nohup podman system service --time=0 "${DOCKER_HOST}" >/tmp/salva-food-podman-service.log 2>&1 &
  for _ in {1..50}; do
    [[ -S "${XDG_RUNTIME_DIR}/podman/podman.sock" ]] && break
    sleep 0.1
  done
fi

if [[ ! -S "${XDG_RUNTIME_DIR}/podman/podman.sock" ]]; then
  cat /tmp/salva-food-podman-service.log 2>/dev/null || true
  exit 1
fi

exec podman compose "$@"

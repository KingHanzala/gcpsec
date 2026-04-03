#!/usr/bin/env sh
set -eu

REPO="KingHanzala/gcpsec"
BINARY_NAME="gcpsec"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

need_cmd uname
need_cmd mktemp
need_cmd tar
need_cmd chmod
need_cmd install

ARCH="$(uname -m)"
case "${ARCH}" in
  x86_64|amd64)
    GOARCH="amd64"
    ;;
  aarch64|arm64)
    GOARCH="arm64"
    ;;
  *)
    echo "Unsupported Linux architecture: ${ARCH}" >&2
    exit 1
    ;;
esac

if command -v curl >/dev/null 2>&1; then
  DOWNLOADER="curl -fsSL"
elif command -v wget >/dev/null 2>&1; then
  DOWNLOADER="wget -qO-"
else
  echo "Missing required downloader: curl or wget" >&2
  exit 1
fi

TMP_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "${TMP_DIR}"
}
trap cleanup EXIT INT TERM

LATEST_URL="https://github.com/${REPO}/releases/latest/download"
ARCHIVE_NAME="${BINARY_NAME}_linux_${GOARCH}.tar.gz"
ARCHIVE_PATH="${TMP_DIR}/${ARCHIVE_NAME}"

echo "Downloading ${ARCHIVE_NAME}..."
sh -c "${DOWNLOADER} \"${LATEST_URL}/${ARCHIVE_NAME}\" > \"${ARCHIVE_PATH}\""

tar -xzf "${ARCHIVE_PATH}" -C "${TMP_DIR}"
chmod +x "${TMP_DIR}/${BINARY_NAME}"
mkdir -p "${INSTALL_DIR}"
install -m 0755 "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"

echo "Installed ${BINARY_NAME} to ${INSTALL_DIR}/${BINARY_NAME}"
echo "Run: ${BINARY_NAME} version"

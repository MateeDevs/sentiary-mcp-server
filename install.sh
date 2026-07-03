#!/bin/sh
set -eu

REPO="MateeDevs/sentiary-mcp-server"
BINARY="sentiary-mcp-server"

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"

case "$os" in
  darwin|linux) ;;
  *)
    echo "Unsupported OS: $os" >&2
    exit 1
    ;;
esac

case "$arch" in
  x86_64|amd64)
    arch="amd64"
    ;;
  arm64|aarch64)
    arch="arm64"
    ;;
  *)
    echo "Unsupported architecture: $arch" >&2
    exit 1
    ;;
esac

if [ "$os" = "linux" ] && [ "$arch" = "arm64" ]; then
  echo "linux/arm64 is not published." >&2
  exit 1
fi

asset="${BINARY}-${os}-${arch}"
archive="${asset}.tar.gz"
url="https://github.com/${REPO}/releases/latest/download/${archive}"
install_dir="${INSTALL_DIR:-}"

if [ -z "$install_dir" ]; then
  if [ -w "/usr/local/bin" ]; then
    install_dir="/usr/local/bin"
  else
    install_dir="$HOME/.local/bin"
  fi
fi

mkdir -p "$install_dir"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT INT TERM

if command -v curl >/dev/null 2>&1; then
  curl -fsSL "$url" -o "$tmp_dir/$archive"
elif command -v wget >/dev/null 2>&1; then
  wget -q "$url" -O "$tmp_dir/$archive"
else
  echo "curl or wget is required." >&2
  exit 1
fi

tar -xzf "$tmp_dir/$archive" -C "$tmp_dir"

if [ -f "$tmp_dir/$asset" ]; then
  source_binary="$tmp_dir/$asset"
elif [ -f "$tmp_dir/$BINARY" ]; then
  source_binary="$tmp_dir/$BINARY"
else
  echo "Archive does not contain $asset or $BINARY." >&2
  exit 1
fi

install_path="$install_dir/$BINARY"
cp "$source_binary" "$install_path"
chmod +x "$install_path"

if command -v "$install_path" >/dev/null 2>&1; then
  version="$($install_path --version 2>/dev/null || true)"
else
  version=""
fi

echo "Installed $BINARY to $install_path"
if [ -n "$version" ]; then
  echo "Version: $version"
fi

case ":$PATH:" in
  *":$install_dir:"*) ;;
  *) echo "Add $install_dir to PATH if your shell cannot find $BINARY." ;;
esac

#!/bin/bash

set -e

APP_NAME="dns-monitor"
VERSION="1.0.0"
INSTALL_DIR="/usr/local/bin"
GITHUB_REPO="mktbsh/dns-monitor"

detect_platform() {
    local os
    local arch
    
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *)          echo "Unsupported OS: $(uname -s)"; exit 1 ;;
    esac
    
    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64" ;;
        arm64|aarch64)  arch="arm64" ;;
        *)              echo "Unsupported architecture: $(uname -m)"; exit 1 ;;
    esac
    
    echo "${os}/${arch}"
}

download_binary() {
    local platform="$1"
    local os="${platform%/*}"
    local arch="${platform#*/}"
    
    local binary_name="${APP_NAME}-${VERSION}-${os}-${arch}"
    local download_url
    local temp_file
    
    if [ "$os" = "windows" ]; then
        binary_name="${binary_name}.exe"
        download_url="https://github.com/${GITHUB_REPO}/releases/download/v${VERSION}/${binary_name}.zip"
        temp_file="/tmp/${binary_name}.zip"
    else
        download_url="https://github.com/${GITHUB_REPO}/releases/download/v${VERSION}/${binary_name}.tar.gz"
        temp_file="/tmp/${binary_name}.tar.gz"
    fi
    
    echo "Downloading ${APP_NAME} v${VERSION} for ${platform}..."
    echo "URL: ${download_url}"
    
    if command -v curl >/dev/null 2>&1; then
        curl -sL "${download_url}" -o "${temp_file}"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "${download_url}" -O "${temp_file}"
    else
        echo "Error: Neither curl nor wget found. Please install one of them."
        exit 1
    fi
    
    echo "${temp_file}"
}

extract_and_install() {
    local temp_file="$1"
    local platform="$2"
    local os="${platform%/*}"
    
    echo "Extracting and installing..."
    
    local temp_dir
    temp_dir="$(mktemp -d)"
    
    if [ "$os" = "windows" ]; then
        if command -v unzip >/dev/null 2>&1; then
            unzip -q "${temp_file}" -d "${temp_dir}"
        else
            echo "Error: unzip not found. Please install unzip."
            exit 1
        fi
    else
        tar -xzf "${temp_file}" -C "${temp_dir}"
    fi
    
    local binary_file
    binary_file="$(find "${temp_dir}" -name "${APP_NAME}*" -type f | head -n1)"
    
    if [ -z "${binary_file}" ]; then
        echo "Error: Binary file not found in archive"
        exit 1
    fi
    
    mkdir -p "${INSTALL_DIR}"
    if [ ! -w "${INSTALL_DIR}" ]; then
        echo "Installing to ${INSTALL_DIR} (requires sudo)..."
        sudo cp "${binary_file}" "${INSTALL_DIR}/${APP_NAME}"
        sudo chmod +x "${INSTALL_DIR}/${APP_NAME}"
    else
        echo "Installing to ${INSTALL_DIR}..."
        cp "${binary_file}" "${INSTALL_DIR}/${APP_NAME}"
        chmod +x "${INSTALL_DIR}/${APP_NAME}"
    fi
    
    rm -rf "${temp_dir}" "${temp_file}"
}

local_install() {
    echo "Local installation mode"
    
    if [ ! -f "main.go" ]; then
        echo "Error: main.go not found. Please run this script from the project root."
        exit 1
    fi
    
    echo "Building from source..."
    go build -ldflags="-X main.Version=${VERSION} -s -w" -o "${APP_NAME}" .
    
    mkdir -p "${INSTALL_DIR}"
    if [ ! -w "${INSTALL_DIR}" ]; then
        echo "Installing to ${INSTALL_DIR} (requires sudo)..."
        sudo cp "${APP_NAME}" "${INSTALL_DIR}/${APP_NAME}"
        sudo chmod +x "${INSTALL_DIR}/${APP_NAME}"
    else
        echo "Installing to ${INSTALL_DIR}..."
        cp "${APP_NAME}" "${INSTALL_DIR}/${APP_NAME}"
        chmod +x "${INSTALL_DIR}/${APP_NAME}"
    fi
    
    rm -f "${APP_NAME}"
}

verify_installation() {
    local binary_path="${INSTALL_DIR}/${APP_NAME}"
    
    if [ -f "${binary_path}" ] && [ -x "${binary_path}" ]; then
        echo "Installation successful!"
        echo "Binary installed at: ${binary_path}"
        echo "Version: $(${binary_path} --version)"
        echo ""
        echo "Usage examples:"
        echo "  ${binary_path} example.com"
        echo "  ${APP_NAME} example.com  # (if ${INSTALL_DIR} is in PATH)"
        echo "  ${APP_NAME} -i 30s example.com api.example.com"
        echo "  ${APP_NAME} -t CNAME --until-change www.example.com"
        echo ""
        echo "Run '${binary_path} --help' for more options."
        
        if ! command -v "${APP_NAME}" >/dev/null 2>&1; then
            echo ""
            echo "Note: ${INSTALL_DIR} is not in your PATH."
            echo "Add the following to your shell profile to use '${APP_NAME}' directly:"
            echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
        fi
    else
        echo "Installation failed. Binary not found at ${binary_path}"
        exit 1
    fi
}

show_help() {
    cat << EOF
DNS Monitor Installation Script

USAGE:
    $0 [OPTIONS]

OPTIONS:
    --local     Install from local source code (requires Go)
    --dir DIR   Set installation directory (default: /usr/local/bin)
    --help      Show this help message

EXAMPLES:
    $0                    # Install latest release from GitHub
    $0 --local            # Build and install from source
    $0 --dir ~/.local/bin # Install to custom directory

REQUIREMENTS:
    - curl or wget (for downloading releases)
    - tar and gzip (for extracting archives)
    - sudo access (if installing to system directories)
    - Go 1.19+ (for --local installation)

EOF
}

main() {
    local local_mode=false
    
    while [ $# -gt 0 ]; do
        case "$1" in
            --local)
                local_mode=true
                shift
                ;;
            --dir)
                INSTALL_DIR="$2"
                shift 2
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            *)
                echo "Unknown option: $1"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done
    
    echo "DNS Monitor Installer v${VERSION}"
    echo "Installing to: ${INSTALL_DIR}"
    echo ""
    
    if [ "${local_mode}" = true ]; then
        local_install
    else
        local platform
        platform="$(detect_platform)"
        echo "Detected platform: ${platform}"
        
        local temp_file
        temp_file="$(download_binary "${platform}")"
        extract_and_install "${temp_file}" "${platform}"
    fi
    
    verify_installation
}

if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    main "$@"
fi
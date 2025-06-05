#!/bin/bash

set -e

VERSION="1.0.0"
APP_NAME="dns-monitor"
BUILD_DIR="build"

echo "Building ${APP_NAME} v${VERSION}"

rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r -a platform_split <<< "$platform"
    GOOS="${platform_split[0]}"
    GOARCH="${platform_split[1]}"
    
    output_name="${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}"
    
    if [ "$GOOS" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    echo "Building for ${GOOS}/${GOARCH}..."
    
    env GOOS="$GOOS" GOARCH="$GOARCH" go build \
        -ldflags="-X main.Version=${VERSION} -s -w" \
        -o "${BUILD_DIR}/${output_name}" \
        .
        
    if [ "$GOOS" != "windows" ]; then
        chmod +x "${BUILD_DIR}/${output_name}"
    fi
done

echo "Build completed. Binaries are in ${BUILD_DIR}/"
ls -la ${BUILD_DIR}/

echo ""
echo "Creating archives..."
cd ${BUILD_DIR}

for file in *; do
    if [ "$file" != "*.tar.gz" ] && [ "$file" != "*.zip" ]; then
        if [[ "$file" == *".exe" ]]; then
            zip "${file%.exe}.zip" "$file"
        else
            tar -czf "${file}.tar.gz" "$file"
        fi
    fi
done

echo "Archives created:"
ls -la *.tar.gz *.zip 2>/dev/null || true

cd ..
echo "Build process completed successfully!"
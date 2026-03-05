#!/bin/bash
# script/update-mockery-config.sh
set -e

CONFIG_FILE=".mockery.yaml"
TMP_FILE="${CONFIG_FILE}.tmp"

echo "🔍 Searching genmock.go files..."

cat > "$TMP_FILE" << 'EOF'
packages:
EOF

# Ищем все genmock.go
find ./internal -type f -name "genmock.go" | while read -r file; do

    dir=$(dirname "$file")
    import_path=$(cd "$dir" && go list -f '{{.ImportPath}}')

    # Берем все интерфейсы без проверки на *Mock
    interfaces=$(grep -E "type .* interface" "$file" | awk '{print $2}')

    [ -z "$interfaces" ] && continue

    echo "  $import_path:" >> "$TMP_FILE"
    echo "    config:" >> "$TMP_FILE"
    echo "      dir: \"{{.InterfaceDir}}/mocks\"" >> "$TMP_FILE"
    echo "      filename: \"{{.InterfaceName | snakecase}}_mock.go\"" >> "$TMP_FILE"
    echo "      pkgname: mocks" >> "$TMP_FILE"
    echo "    interfaces:" >> "$TMP_FILE"

    for iface in $interfaces; do
        echo "      $iface: {}" >> "$TMP_FILE"
    done

done

mv "$TMP_FILE" "$CONFIG_FILE"

echo "✅ Generated $CONFIG_FILE"
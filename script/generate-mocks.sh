#!/bin/bash
# script/generate-mocks.sh
set -e

MOCKERY="./bin/mockery"

if [ ! -f "$MOCKERY" ]; then
    echo "❌ mockery not found"
    exit 1
fi

echo "🔧 Using mockery"

"$MOCKERY"
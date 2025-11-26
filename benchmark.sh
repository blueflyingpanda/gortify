#!/bin/bash

# Install hey if not present: go install github.com/rakyll/hey@latest

API_URL="http://localhost:8080"
REQUESTS=500
WORKERS=5

codes=(
    "60cf0877"
)

benchmark() {
    local name=$1
    echo "=========================================="
    echo "Testing: $name"
    echo "=========================================="

    echo "Combined test (random codes):"
    # Run mixed test
    /Users/ilya/go/bin/hey -n $REQUESTS -c $WORKERS -disable-redirects "${API_URL}/${codes[0]}" | tee "${name}_summary.txt"
    echo ""
}

echo "Step 1: WITHOUT Redis"
echo "Press Enter to start..."
read
benchmark "without_redis"

echo "Step 2: Enable Redis and restart server"
echo "Press Enter when ready..."
read

benchmark "with_redis"

echo "Done! Check *_summary.txt files"

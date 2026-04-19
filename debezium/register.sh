#!/bin/sh
set -e

CONNECT_URL="${CONNECT_URL:-http://localhost:8083}"
CONNECTOR_NAME="outbox-connector"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "Waiting for Kafka Connect to be ready..."
until curl -sf "${CONNECT_URL}/connectors" > /dev/null; do
  sleep 2
done

echo "Checking connector '${CONNECTOR_NAME}'..."
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "${CONNECT_URL}/connectors/${CONNECTOR_NAME}")

if [ "$HTTP_STATUS" = "200" ]; then
  echo "Connector already exists. To update it, delete first:"
  echo "  curl -X DELETE ${CONNECT_URL}/connectors/${CONNECTOR_NAME}"
  exit 0
fi

echo "Registering connector..."
curl -sf -X POST "${CONNECT_URL}/connectors" \
  -H "Content-Type: application/json" \
  -d @"${SCRIPT_DIR}/connector.json"

echo ""
echo "Connector registered successfully."

#!/bin/bash

# Display error and exit
function error_exit() {
  echo "Error: $1" >&2
  exit 1
}

# Parse parameters
while [[ "$#" -gt 0 ]]; do
  case $1 in
    --path)
      TARGET_PATH="$2"
      shift 2
      ;;
    *)
      error_exit "Unknown parameters $1"
      ;;
  esac
done

# Check --path
if [[ -z "$TARGET_PATH" ]]; then
  error_exit "Please input the --path parameter."
fi


if [[ "$TARGET_PATH" = /* ]]; then
  ABSOLUTE_PATH="$TARGET_PATH"
else
  ABSOLUTE_PATH=$(cd "$TARGET_PATH" 2>/dev/null && pwd) || error_exit "Cannot resolve path $TARGET_PATH"
fi

if [[ ! -d "$ABSOLUTE_PATH" ]]; then
  mkdir -p "$ABSOLUTE_PATH" || error_exit "Cannot create path $ABSOLUTE_PATH"
fi

docker run --rm -it -v "$ABSOLUTE_PATH:/app" -w /app krakend/builder:2.7.2 go build -buildmode=plugin -o redis-rate-limit.so .
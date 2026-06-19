#!/bin/bash

tfvars_file="terraform.tfvars"
dry_run=false
webhooks=()
channels=()
severity=()
block_starts=()
block_ends=()

current_webhook=""
current_channel=""
current_severity=""
current_block_start=""

usage() {
  cat <<'EOF'
Usage: bash slack_webhook_test.bash [--dry-run]

Options:
  --dry-run  Print the blocks that would be uncommented without editing terraform.tfvars.
  --help     Show this help text.
EOF
}

uncomment_block() {
  local file="$1"
  local start_line="$2"
  local end_line="$3"

  perl -i -pe 'if ($. >= '"$start_line"' && $. <= '"$end_line"') { s/^(\s*)#\s?/$1/ }' "$file"
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run)
      dry_run=true
      ;;
    --help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
  shift
done

line_number=0
while IFS= read -r line; do
  ((line_number++))

  if [[ "$line" =~ ^[[:space:]]*#[[:space:]]*\{ ]]; then
    current_block_start="$line_number"
  fi

  if [[ "$line" =~ ^[[:space:]]*#[[:space:]]*webhook[[:space:]]*=[[:space:]]*\"([^\"]+)\" ]]; then
    current_webhook="${BASH_REMATCH[1]}"
  elif [[ "$line" =~ ^[[:space:]]*#[[:space:]]*channel[[:space:]]*=[[:space:]]*\"([^\"]+)\" ]]; then
    current_channel="${BASH_REMATCH[1]}"
  elif [[ "$line" =~ ^[[:space:]]*#[[:space:]]*severity[[:space:]]*=[[:space:]]*\"([^\"]+)\" ]]; then
    current_severity="${BASH_REMATCH[1]}"
  fi

  if [[ "$line" =~ ^[[:space:]]*#[[:space:]]*\},?[[:space:]]*$ ]] && [[ -n "$current_webhook" && -n "$current_channel" && -n "$current_severity" && -n "$current_block_start" ]]; then
    webhooks+=("$current_webhook")
    channels+=("$current_channel")
    severity+=("$current_severity")
    block_starts+=("$current_block_start")
    block_ends+=("$line_number")
    current_webhook=""
    current_channel=""
    current_severity=""
    current_block_start=""
  fi
done < "$tfvars_file"

for i in "${!webhooks[@]}"; do
  webhook="${webhooks[$i]}"
  channel="${channels[$i]}"
  sev="${severity[$i]}"
  block_start="${block_starts[$i]}"
  block_end="${block_ends[$i]}"
  echo -e "\n---"
  echo -e "webhook: $webhook"
  echo "channel: $channel"
  echo "severity: $sev"

  if [[ "$dry_run" == true ]]; then
    echo "Dry run: would test webhook and uncomment lines $block_start-$block_end in $tfvars_file if the response is ok"
    echo -e "\n---\n"
    continue
  fi

  response=$(curl -sS -X POST -H 'Content-type: application/json' \
    --data "{\"text\":\"CP Alert manager test for channel $channel with severity $sev\"}" \
    "$webhook" || true)

  echo "$response"

  if [[ "$response" == "ok" ]]; then
    echo "Detected ok response for $channel. Uncommenting lines $block_start-$block_end in $tfvars_file"
    uncomment_block "$tfvars_file" "$block_start" "$block_end"
  fi

   echo -e "---\n"
done
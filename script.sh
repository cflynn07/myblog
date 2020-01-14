#!/bin/bash

overrides=(
  "develop_image=$HOME"
  "deployment_sha=bang"
  "action_id=fly"
)

echo "0: ${overrides[*]}"
overrides=$(for i in "${overrides[@]}"; do echo -n "$i,"; done)
echo "1: $overrides"
echo "2: ${#overrides}"
echo "3: ${overrides:0:${#overrides}-1}"

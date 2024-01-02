#! /bin/bash

. .env

curl \
  -H "accept: application/json"\
  -H "authorization: Bearer $BERYX_TOKEN" \
  "https://api.zondax.ch/fil/data/v3/mainnet/transactions/address/f410f6dy45thxrvar53m4ugimu7yvofamzmwtxrc4aaq/receiver?limit=100&remove_internal_txs=1" | jq


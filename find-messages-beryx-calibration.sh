#! /bin/bash

. .env

# Agent address

#curl \
#  -H "accept: application/json"\
#  -H "authorization: Bearer $BERYX_TOKEN" \
#  "https://api.zondax.ch/fil/data/v3/calibration/transactions/address/t410fykggyf5ye42oxqouavikwguici5hupvq4c3fuxi/receiver?limit=100" | jq

# Receive tx
curl \
  -H "accept: application/json"\
  -H "authorization: Bearer $BERYX_TOKEN" \
  https://api.zondax.ch/fil/data/v3/calibration/transactions/id/MTIzMzIyMy9iYWZ5MmJ6YWNlYm5qdGNuYmxlaXg3d2xzZnZ4aWp4bTVlcmRxNWFrZHpsYXB5d2o3NHhtNHEzdnNzanZuNi9iYWZ5MmJ6YWNlY2M2ajVibGFnbnF5ZXMzZm01cHJldXV2NG5wY3NkY2JyZXh3enR1YjZ1Y2hneTNnbmltdy9iYWZ5MmJ6YWNlYXFwajZqdXV2b3ljZzd4eGNtNGtyb2dveXhxZWc2dmMycTc1Mm9yZjRkcDZzdjd3dnRiay9lNGM2ZDM5MC0xODUwLTVjNzYtYjQ3MC1kNWQ2ZGQ3YTQ1YWQ= | jq

# Owner address

#curl \
#  -H "accept: application/json"\
#  -H "authorization: Bearer $BERYX_TOKEN" \
#  "https://api.zondax.ch/fil/data/v3/calibration/transactions/address/t045541/receiver?limit=100&remove_internal_txs=1" | jq



export SWITCHBOT_TOKEN=$(vault kv get -format=json kv/switchbot|jq -r '.data.data.TOKEN')
export SWITCHBOT_CLIENT_SECRET=$(vault kv get -format=json kv/switchbot|jq -r '.data.data.CLIENT_SECRET')

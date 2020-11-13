#!/bin/bash

# Deploys a vault server with a configured PKI backend
#  to the host system using docker

set -e

cd "$(dirname "$0")"

if [ -z ${LOCAL_VAULT_CONTAINER_NAME+x} ]; then echo "LOCAL_VAULT_CONTAINER_NAME is unset"; exit 1; fi
if [ -z ${LOCAL_VAULT_TOKEN+x} ]; then echo "LOCAL_VAULT_TOKEN is unset"; exit 1; fi
if [ -z ${LOCAL_VAULT_PORT+x} ]; then echo "LOCAL_VAULT_PORT is unset"; exit 1; fi
if [ -z ${LOCAL_VAULT_ADDR+x} ]; then echo "LOCAL_VAULT_ADDR is unset"; exit 1; fi
if [ -z ${LOCAL_VAULT_PKI_URL+x} ]; then echo "LOCAL_VAULT_PKI_URL is unset"; exit 1; fi
if [ -z ${LOCAL_VAULT_CERT_URL+x} ]; then echo "LOCAL_VAULT_CERT_URL is unset"; exit 1; fi
if [ -z ${LOCAL_VAULT_BACKEND_PATH+x} ]; then echo "LOCAL_VAULT_BACKEND_PATH is unset"; exit 1; fi
if [ -z ${LOCAL_VAULT_BACKEND_NAME+x} ]; then echo "LOCAL_VAULT_BACKEND_NAME is unset"; exit 1; fi
if [ -z ${LOCAL_VAULT_BACKEND_DOMAIN+x} ]; then echo "LOCAL_VAULT_BACKEND_DOMAIN is unset"; exit 1; fi

clean() {
    docker ps | grep $LOCAL_VAULT_CONTAINER_NAME | \
        awk '{print $1}' | \
        xargs -I{} docker rm {} --force >> /dev/null
    rm -f terraform.tfstate*
}
trap clean ERR

# run a local vault instance in dev mode
# https://registry.hub.docker.com/_/vault/
vault_server() {
    echo "Starting development vault instance on port \"${LOCAL_VAULT_PORT}\" with root token \"${LOCAL_VAULT_TOKEN}\""
    docker run -d \
        --cap-add=IPC_LOCK \
        --name=$LOCAL_VAULT_CONTAINER_NAME \
        -p $LOCAL_VAULT_PORT:8200 \
        -e "VAULT_DEV_ROOT_TOKEN_ID=${LOCAL_VAULT_TOKEN}" \
        -e "VAULT_TOKEN=${LOCAL_VAULT_TOKEN}" \
        -e "VAULT_UI=true" \
        -e "VAULT_ADDR=http://localhost:8200" \
        vault:1.6.0 >> /dev/null

    sleep 5

    docker exec $LOCAL_VAULT_CONTAINER_NAME \
        vault secrets enable pki >> /dev/null
    docker exec $LOCAL_VAULT_CONTAINER_NAME \
        vault secrets tune -max-lease-ttl=8801h pki
    docker exec $LOCAL_VAULT_CONTAINER_NAME \
        vault write pki/root/generate/internal common_name=$LOCAL_VAULT_BACKEND_DOMAIN ttl=8801h >> /dev/null
    docker exec $LOCAL_VAULT_CONTAINER_NAME \
        vault write pki/roles/$LOCAL_VAULT_BACKEND_NAME allowed_domains=$LOCAL_VAULT_BACKEND_DOMAIN allow_subdomains=true max_ttl=8800h >> /dev/null
}

clean
vault_server
echo "Vault can be reached at ${LOCAL_VAULT_ADDR}/ui, use token \"${LOCAL_VAULT_TOKEN}\""

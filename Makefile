ifeq (, $(shell which docker))
    $(error "No docker in $(PATH)")
endif

VERSION := $(shell cat cmd/version.go | grep "const VERSION" | cut -c 17- | tr -d '"')

# for tests
export LOCAL_VAULT_CONTAINER_NAME := dev-vault
export LOCAL_VAULT_TOKEN := token
export LOCAL_VAULT_PORT := 58200
export LOCAL_VAULT_ADDR := http://localhost:58200
export LOCAL_VAULT_BACKEND_PATH := pki
export LOCAL_VAULT_BACKEND_NAME := test-dot-local
export LOCAL_VAULT_BACKEND_DOMAIN := test.local # test/terraform/pki.tf assumes the use of this domain
export LOCAL_VAULT_PKI_URL := http://localhost:58200/v1/pki
export LOCAL_VAULT_CERT_URL := http://localhost:58200/v1/pki/issue/test-dot-local

$(shell mkdir -p artifacts)

build: clean local-vault
	$(info building bmcert ${VERSION})

	@docker build \
		--network=host \
		--no-cache \
		--build-arg version=${VERSION} \
	    --build-arg token=${LOCAL_VAULT_TOKEN} \
	    --build-arg addr=${LOCAL_VAULT_ADDR} \
	    --build-arg url=${LOCAL_VAULT_CERT_URL} \
		--build-arg pki_url=${LOCAL_VAULT_PKI_URL} \
	    -t bmcert:${VERSION} .

	@docker create -ti --name bmcertartifacts bmcert:${VERSION} bash && \
	    docker cp bmcertartifacts:/bmcert/bmcert-v${VERSION}-linux-amd64.zip artifacts/bmcert-v${VERSION}-linux-amd64.zip && \
	    docker cp bmcertartifacts:/bmcert/bmcert-v${VERSION}-darwin-amd64.zip artifacts/bmcert-v${VERSION}-darwin-amd64.zip && \
	    docker cp bmcertartifacts:/bmcert/bmcert-v${VERSION}.SHA256SUMS artifacts/bmcert-v${VERSION}.SHA256SUMS

	# cleanup
	@docker rm -fv bmcertartifacts &> /dev/null

local-vault:
	scripts/local_vault.sh

test:
	go test ./...

lint:
	golint ./...

clean:
	$(shell rm -rf artifacts/*)
	$(shell rm -rf test/terraform.tfstate*)
	$(shell docker ps -a | grep 'bmcertartifacts' | awk '{print $1}' | xargs -n1 docker rm)

ifeq (, $(shell which docker))
    $(error "No docker in $(PATH)")
endif

ifndef VAULT_ADDR
    $(error VAULT_ADDR is not set.)
endif


ifndef VAULT_CERT_URL
    $(error VAULT_CERT_URL is not set)
endif

ifndef VAULT_PKI_URL
    $(error VAULT_PKI_URL is not set)
endif

VERSION := $(shell cat cmd/version.go | grep "const VERSION" | cut -c 17- | tr -d '"')

$(shell mkdir -p artifacts)

build: clean
	$(info building bmcert ${VERSION})

	@docker build \
		--no-cache \
		--build-arg version=${VERSION} \
	    --build-arg token=${VAULT_GITHUB_TOKEN} \
	    --build-arg addr=${VAULT_ADDR} \
	    --build-arg url=${VAULT_CERT_URL} \
		--build-arg pki_url=${VAULT_PKI_URL} \
	    -t bmcert:${VERSION} .

	@docker create -ti --name bmcertartifacts bmcert:${VERSION} bash && \
	    docker cp bmcertartifacts:/bmcert/bmcert-v${VERSION}-linux-amd64.zip artifacts/bmcert-v${VERSION}-linux-amd64.zip && \
	    docker cp bmcertartifacts:/bmcert/bmcert-v${VERSION}-darwin-amd64.zip artifacts/bmcert-v${VERSION}-darwin-amd64.zip && \
	    docker cp bmcertartifacts:/bmcert/bmcert-v${VERSION}.SHA256SUMS artifacts/bmcert-v${VERSION}.SHA256SUMS

	# cleanup
	@docker rm -fv bmcertartifacts &> /dev/null

test:
	go test ./...

lint:
	golint ./...

clean:
	$(shell rm -rf artifacts/*)
	$(shell docker ps -a | grep 'bmcertartifacts' | awk '{print $1}' | xargs -n1 docker rm)

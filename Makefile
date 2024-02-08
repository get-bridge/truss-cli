VERSION = v0.3.6

release:
	git tag -a ${VERSION} -m ${VERSION}
	git push origin ${VERSION}

dryrun:
	goreleaser --snapshot --skip-publish --rm-dist

test:
	go mod download
	go generate ./...
	go test ./cmd/ ./truss/ -timeout 15000ms

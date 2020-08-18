VERSION = v0.0.10

release:
	git tag -a ${VERSION} -m ${VERSION}
	git push origin ${VERSION}

dryrun:
	goreleaser --snapshot --skip-publish --rm-dist

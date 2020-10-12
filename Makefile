VERSION = v0.1.1

release:
	git tag -a ${VERSION} -m ${VERSION}
	git push origin ${VERSION}

dryrun:
	goreleaser --snapshot --skip-publish --rm-dist

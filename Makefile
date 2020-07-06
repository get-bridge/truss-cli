VERSION = v0.0.8

release:
	git tag -a ${VERSION} -m ${VERSION}
	git push origin ${VERSION}
	goreleaser release --rm-dist

dryrun:
	goreleaser --snapshot --skip-publish --rm-dist

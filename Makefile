# Generate manifests for CRDs

.PHONY: manifests

CONTROLLER_GEN = /home/sahadat/go/bin/controller-gen

manifests:
	$(CONTROLLER_GEN) crd:trivialVersions=true paths="./..." output:crd:artifacts:config=config/crd/bases
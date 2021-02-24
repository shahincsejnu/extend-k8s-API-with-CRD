# Generate manifests for CRDs

.PHONY: manifests

CONTROLLER_GEN = ./bin/controller-gen

manifests:
	@$(CONTROLLER_GEN) crd:trivialVersions=true paths="./pkg/..." output:crd:artifacts:config=config/crds


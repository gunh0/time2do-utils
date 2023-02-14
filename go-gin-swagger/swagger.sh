export PATH=$(go env GOPATH)/bin:$PATH
swag init --parseDependency --parseInternal --parseDepth 1 -md ./documentation
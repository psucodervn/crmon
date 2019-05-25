build:
	 go build \
	 	-ldflags "-X 'main.version=$(git describe --tags --abbrev=0)' -X 'main.build=$(git rev-parse --short HEAD)'" \
	 	-o crman cmd/crman/*

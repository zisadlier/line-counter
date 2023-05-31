GO := go
GOBUILD := $(GO) build
GOTEST := $(GO) test
GOMOD := $(GO) mod
GOFLAGS := -v

build:
	$(GOBUILD) $(GOFLAGS) -o "linecounter" cmd/main.go

test:
	$(GOTEST) $(GOFLAGS) ./...

app_test:
	cd test_files; bash application_tests.sh; cd ..

tidy:
	$(GOMOD) tidy
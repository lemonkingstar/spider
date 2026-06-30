LOCAL_PATH = $(shell pwd)

.PHONY: tagger
# tagger
tagger:
	protoc --proto_path=./ tagger/tagger.proto --go_out=. --go_opt=paths=source_relative

.PHONY: example
# example
example:
	protoc --proto_path=./ example/example.proto --go_out=. --go_opt=paths=source_relative
	protoc --proto_path=./ example/example.proto --gotags_out=. --gotags_opt=paths=source_relative
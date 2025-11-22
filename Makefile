
.PHONY: proto
# generate proto files
proto:
	protoc pkg/ptype/*.proto --go_out=. --go_opt=paths=source_relative

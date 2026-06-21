proto:
	protoc --proto_path=./ --proto_path=../ example.proto --go_out=. --go_opt=paths=source_relative
	protoc --proto_path=./ --proto_path=../ example.proto --gotags_out=. --gotags_opt=paths=source_relative

.PHONY: proto
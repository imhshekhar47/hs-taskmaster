EXE_NAME=skill-api

clean:
	if [ -d "./build" ]; then rm -r "./build"; fi;
	if [ -d "./pb" ]; then find ./pb -type f -exec rm {} \; ; fi;
	

generate: clean
	if [ ! -d "./pb" ]; then mkdir -p "./pb"; fi;
	if [ ! -d "./build/openapi/v2" ]; then mkdir -p "./build/openapi/v2"; fi;

	protoc -I="proto" \
	--go_out="pb" --go_opt=paths=source_relative \
	--go-grpc_out="pb" --go-grpc_opt=paths=source_relative \
	--gorm_out=. \
	--grpc-gateway_out="pb" --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out ./build/openapi/v2 \
	--openapiv2_opt logtostderr=true \
	proto/*.proto

swagger: generate
	find cmd -name *.swagger.json -exec rm {} \;
	cp build/openapi/v2/${EXE_NAME}.swagger.json cmd/

test: swagger
	if [ ! -d "./build/tests" ]; then mkdir -p "./build/tests"; fi;
	gotestsum --format testname \
			--junitfile "./build/tests/unit-tests.xml" \
			-- -coverprofile=build/coverage.out ./...

build: test
	if [ ! -d "./build/bin" ];  then mkdir -p "./build/bin"; fi;
	go build -o ./build/bin/${EXE_NAME}

run: build
	./build/bin/${EXE_NAME}

image:
	docker build -t "${EXE_NAME}:v1" .

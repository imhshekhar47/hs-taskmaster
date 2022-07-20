test: 
	if [ ! -d "./build/tests" ]; then mkdir -p "./build/tests"; fi;
	gotestsum --format testname \
			--junitfile "./build/tests/unit-tests.xml" \
			-- -coverprofile=build/coverage.out ./...
build:
	go build 
SRCS=$(wildcard *.go)

gc: $(SRCS)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ $^

test: gc
	docker run --rm -v $(shell pwd):/gc -w /gc compilerbook bash -c \
		'./test.sh'

clean:
	rm -f gc *.o *~ tmp*

.PHONY: test clean

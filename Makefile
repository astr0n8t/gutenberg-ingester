test:
	go test ./...
fmt:
	go fmt ./...
mem-profile:
	while true; do top -l 1 | grep "gutenberg-ingest" | awk '{print "MEM="$$8}'; sleep .5s; done

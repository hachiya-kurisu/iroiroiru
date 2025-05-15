all: iroiroserve iroiroload

again: clean all

iroiroserve: iroiroiru.go cmd/iroiroserve/main.go
	go build -C cmd/iroiroserve -o ../../iroiroserve

iroiroload: iroiroiru.go cmd/iroiroload/main.go
	go build -C cmd/iroiroload -o ../../iroiroload

clean:
	rm -f iroiroserve iroiroload

test:
	go test -cover

push:
	got send
	git push github

fmt:
	gofmt -s -w *.go cmd/*/main.go

cover:
	go test -coverprofile=cover.out
	go tool cover -html cover.out

README.md: README.gmi
	sisyphus -f markdown <README.gmi >README.md

doc: README.md

publish:
	rsync -raz --safe-links --progress iroiroview/* iroiroiru.jp:/var/iroiro/www

release: push
	git push github --tags

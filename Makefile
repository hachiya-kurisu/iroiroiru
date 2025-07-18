all: iroiroserve iroiroload iroiroscrape

again: clean all

iroiroscrape: iroiroiru.go cmd/iroiroscrape/main.go
	go build -C cmd/iroiroscrape -o ../../iroiroscrape

iroiroserve: iroiroiru.go cmd/iroiroserve/main.go
	go build -C cmd/iroiroserve -o ../../iroiroserve

iroiroload: iroiroiru.go cmd/iroiroload/main.go
	go build -C cmd/iroiroload -o ../../iroiroload

clean:
	rm -f iroiroserve iroiroload iroiroscrape

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

peek:
	caddy file-server --root iroiroview

README.md: README.gmi
	sisyphus -f markdown <README.gmi >README.md

doc: README.md

publish:
	rsync -raz --safe-links --progress iroiroview/* iroiroiru:/var/iroiro/www

release: push
	git push github --tags

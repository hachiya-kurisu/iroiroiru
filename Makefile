all: iroiroiru

again: clean all

iroiroiru: iroiroiru.go cmd/iroiroiru/main.go
	go build -C cmd/iroiroiru -o ../../iroiroiru

clean:
	rm -f iroiroiru

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
	rsync -raz --safe-links --progress www/* iroiroiru.jp:/var/iroiro/www

release: push
	git push github --tags

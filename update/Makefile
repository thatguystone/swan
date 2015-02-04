.PHONY: update update_stopwords

update: update_stopwords

update_stopwords: python-goose
	go build
	./update
	go fmt ../...

clean:
	rm -rf python-goose
	rm -rf python-goose.tar.gz
	go clean

python-goose.tar.gz:
	curl -L https://api.github.com/repos/grangier/python-goose/tarball > python-goose.tar.gz

python-goose: python-goose.tar.gz
	mkdir python-goose
	tar xvf $< --strip-components 1 --directory $@

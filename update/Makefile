.PHONY: update update_tests update_stopwords

update: update_stopwords update_tests

update_stopwords: python-goose
	go build
	./update
	go fmt ../...

update_tests: python-goose
	rm -rf ../test_data/$<
	cp -r $</tests/data ../test_data/$<

clean:
	rm -rf python-goose
	rm -rf python-goose.tar.gz

python-goose.tar.gz:
	curl -L https://api.github.com/repos/grangier/python-goose/tarball > python-goose.tar.gz

python-goose: python-goose.tar.gz
	mkdir python-goose
	tar xvf $< --strip-components 1 --directory $@

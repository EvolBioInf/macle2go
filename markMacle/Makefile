version = $(shell bash ../scripts/version.sh)
date    = $(shell bash ../scripts/date.sh)

exe = markMacle
all: tangle weave

tangle: $(exe)
$(exe): $(exe).go
	go build -ldflags "-X main.version=$(version) -X main.date=$(date)" $(exe).go
$(exe).go: $(exe).org
	bash ../scripts/org2nw $(exe).org | notangle -R$(exe).go | gofmt > $(exe).go
test: $(exe) $(exe)_test.go
	go test -v
$(exe)_test.go: $(exe).org
	bash ../scripts/org2nw $(exe).org | notangle -R$(exe)_test.go | gofmt > $(exe)_test.go

weave: $(exe)Doc.pdf
$(exe)Doc.pdf: $(exe).tex
	latex $(exe)Doc.tex
	bibtex $(exe)Doc
	latex $(exe)Doc.tex
	latex $(exe)Doc.tex
	dvips $(exe)Doc -o -q
	ps2pdf $(exe)Doc.ps
$(exe).tex: $(exe).org
	bash ../scripts/org2nw $(exe).org | sed 's/_test/\\_test/' | noweave -n -x > $(exe).tex

clean:
	rm -f $(exe) $(exe).go $(exe)_test.go $(exe)Doc.ps $(exe)Doc.aux *.bbl *.blg $(exe)Doc.log $(exe)Doc.dvi $(exe)Doc.pdf $(exe).tex

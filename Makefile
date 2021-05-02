TARG := shorten_path
LDFLAGS := "-s -w"

all: \
	$(TARG).linux-arm \
	$(TARG).linux-amd64 \
	$(TARG).freebsd-amd64 \
	$(TARG).openbsd-amd64 \
	$(TARG).netbsd-amd64 \
	$(TARG).darwin-amd64

clean:
	rm -f $(TARG) $(TARG).*-*

$(TARG).linux-amd64: /dev/null
	go build -ldflags $(LDFLAGS)
	mv $(TARG) $@
	upx $@

$(TARG).linux-%: /dev/null
	GOOS=linux GOARCH=$(patsubst $(TARG).linux-%,%,$@) GOARM=7 go build -ldflags $(LDFLAGS)
	mv $(TARG) $@

$(TARG).%-amd64: /dev/null
	GOOS=$(patsubst $(TARG).%-amd64,%,$@) GOARCH=amd64 go build -ldflags $(LDFLAGS)
	mv $(TARG) $@

.PHONY: all clean

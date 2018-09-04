include .env

.PHONY: vendor example refresh-fixtures clean-fixtures

LINODE_FIXTURE_INSTA***REMOVED***E:=76859403
LINODE_FIXTURE_VOLUME:=6574839201

.PHONY: test
test: vendor
	@LINODE_TEST_INSTA***REMOVED***E=$(LINODE_FIXTURE_INSTA***REMOVED***E) \
	LINODE_TEST_VOLUME=$(LINODE_FIXTURE_VOLUME) \
	LINODE_FIXTURE_MODE="play" \
	LINODE_TOKEN="awesometokenawesometokenawesometoken" \
	vgo test $(ARGS)

$(GOPATH)/bin/vgo:
	@go get -u golang.org/x/vgo

vendor: $(GOPATH)/bin/vgo
	@vgo get -u

example:
	@go run example/main.go

clean-fixtures:
	@-rm fixtures/*.yaml

refresh-fixtures: clean-fixtures fixtures

.PHONY: fixtures
fixtures:
	@echo "* Running fixtures"
	@LINODE_TOKEN=$(LINODE_TOKEN) \
	LINODE_TEST_INSTA***REMOVED***E=$(LINODE_TEST_INSTA***REMOVED***E) \
	LINODE_TEST_VOLUME=$(LINODE_TEST_VOLUME) \
	LINODE_FIXTURE_MODE="record" go test $(ARGS)
	@echo "* Santizing fixtures"
	@for yaml in fixtures/*yaml; do \
		sed -E -i "" -e "s/$(LINODE_TOKEN)/awesometokenawesometokenawesometoken/g" \
			-e "s/$(LINODE_TEST_INSTA***REMOVED***E)/$(LINODE_FIXTURE_INSTA***REMOVED***E)/g" \
			-e 's/20[0-9]{2}-[01][0-9]-[0-3][0-9]T[0-2][0-9]:[0-9]{2}:[0-9]{2}/2018-01-02T03:04:05/g' \
			-e 's/nb-[0-9]{1,3}-[0-9]{1,3}-[0-9]{1,3}-[0-9]{1,3}\./nb-10-20-30-40./g' \
			-e 's/192\.168\.((1?[0-9][0-9]?|2[0-4][0-9]|25[0-5])\.)(1?[0-9][0-9]?|2[0-4][0-9]|25[0-5])/192.168.030.040/g' \
			-e '/^192\.168/!s/((1?[0-9][0-9]?|2[0-4][0-9]|25[0-5])\.){3}(1?[0-9][0-9]?|2[0-4][0-9]|25[0-5])/010.020.030.040/g' \
			-e 's/(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))/1234::5678/g' \
			-e "s/$(LINODE_TEST_VOLUME)/$(LINODE_FIXTURE_VOLUME)/g" $$yaml; \
	done

.PHONY: godoc
godoc:
	@godoc -http=:6060

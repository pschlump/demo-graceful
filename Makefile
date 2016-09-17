
all:
	go build

test: testSetup test1 test2

testSetup:
	@rm -f out/test*.out

test1:
	@echo "Send Request, verify get success and 200 status code"
	@mkdir -p ./out
	@curl -i -X GET 'http://localhost:8123/api/status' >out/test1.out 2>/dev/null
	@grep '^HTTP.*200 OK' out/test1.out >/dev/null
	@grep '{"status":"success"}' out/test1.out >/dev/null


test2:
	@echo "Send Request, verify get success"

xx:
	@mkdir -p ./out
	@curl -i -X GET 'http://localhost:8123/api/status' >out/test1.out
	@grep '^HTTP.*200 OK' out/test1.out >/dev/null
	@grep '{"status":"success"}' out/test1.out >/dev/null



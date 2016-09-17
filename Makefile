
all:
	go build

test: testSetup test1 test4 test2 test3
	@echo PASS

testSetup:
	@go build
	@rm -f out/test*.out

test1:
	@echo "Send Request, verify get success and 200 status code"
	@mkdir -p ./out
	./cleanup-this-one.sh
	@curl -i -X GET 'http://localhost:8123/api/status' >out/test1.out 2>/dev/null
	@grep '^HTTP.*200 OK' out/test1.out >/dev/null
	@grep '{"status":"success"}' out/test1.out >/dev/null


test2:
	@echo "Send Request, verify get success, then shutdown, then verify shutdown pending state"
	@mkdir -p ./out
	curl -i -X GET 'http://localhost:8123/api/status' >out/test2.a.out 2>/dev/null
	grep '^HTTP.*200 OK' out/test2.a.out >/dev/null
	curl -i -X GET 'http://localhost:8123/api/shutdown' >out/test2.b.out 2>/dev/null
	grep '^HTTP.*200 OK' out/test2.b.out >/dev/null
	grep '"status":"success"' out/test2.b.out >/dev/null
	-curl -i -X GET 'http://localhost:8123/api/status' >out/test2.c.out 2>/dev/null

test3:
	@echo "Start server and run tests, multiple request finish after shutdown request"
	@echo "Stop existing server, if any"
	-@curl -i -X GET 'http://localhost:8123/api/shutdown' >/dev/null 2>/dev/null
	@rm -f out/test3*.out
	./cleanup-this-one.sh
	-curl -i -X GET 'http://localhost:8123/api/status' >out/test3.a.out 2>/dev/null
	-curl -i -X GET 'http://localhost:8123/?password=am1' >out/test3.b.out 2>/dev/null &
	-curl -i -X GET 'http://localhost:8123/?password=am2' >out/test3.c.out 2>/dev/null &
	-curl -i -X GET 'http://localhost:8123/?password=am3' >out/test3.d.out 2>/dev/null &
	-curl -i -X GET 'http://localhost:8123/?password=am4' >out/test3.e.out 2>/dev/null &
	curl -i -X GET 'http://localhost:8123/api/shutdown' >out/test3.f.out 2>/dev/null 
	-curl -i -X GET 'http://localhost:8123/?password=am5' >out/test3.g.out 2>out/test3.g.err &
	-curl -i -X GET 'http://localhost:8123/?password=am6' >out/test3.h.out 2>out/test3.h.err &
	-curl -i -X GET 'http://localhost:8123/?password=am7' >out/test3.i.out 2>out/test3.i.err &
	-curl -i -X GET 'http://localhost:8123/?password=am8' >out/test3.j.out 2>out/test3.j.err &
	@sleep 7		# give it plenty of time
	grep '^HTTP.*200 OK' out/test3.a.out >/dev/null
	grep '^HTTP.*200 OK' out/test3.b.out >/dev/null
	grep '^HTTP.*200 OK' out/test3.c.out >/dev/null
	grep '^HTTP.*200 OK' out/test3.d.out >/dev/null
	grep '^HTTP.*200 OK' out/test3.e.out >/dev/null
	grep '^HTTP.*200 OK' out/test3.f.out >/dev/null
	grep 'Connection refused' out/test3.g.err >/dev/null
	grep 'Connection refused' out/test3.h.err >/dev/null
	grep 'Connection refused' out/test3.i.err >/dev/null
	grep 'Connection refused' out/test3.j.err >/dev/null

test4:
	@rm -f out/test4*.out
	./cleanup-this-one.sh
	curl -i -X GET 'http://localhost:8123/?password=angryMonkey' >out/test4.a.out 2>/dev/null 
	tail -1 out/test4.a.out >out/test4.b.out
	diff -w ref/test4.a.out out/test4.b.out

test5:
	@rm -f out/test4*.out
	./cleanup-this-one.sh
	kill -HUP `cat pid-of-svr`
	./test5.sh

	

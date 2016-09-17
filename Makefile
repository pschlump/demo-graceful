
all:
	go build

test: testSetup test1 test4 test2 test3 test5
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
	@echo PASS


test2:
	@echo "Send Request, verify get success, then shutdown, then verify shutdown pending state"
	@mkdir -p ./out
	curl -i -X GET 'http://localhost:8123/api/status' >out/test2.a.out 2>/dev/null
	grep '^HTTP.*200 OK' out/test2.a.out >/dev/null
	curl -i -X GET 'http://localhost:8123/api/shutdown' >out/test2.b.out 2>/dev/null
	grep '^HTTP.*200 OK' out/test2.b.out >/dev/null
	grep '"status":"success"' out/test2.b.out >/dev/null
	@echo "curl should report an error - since the server should not be running"
	-curl -i -X GET 'http://localhost:8123/api/status' >out/test2.c.out 2>/dev/null
	@echo PASS

test3:
	@echo "Start server and run tests, multiple request finish after shutdown request"
	@echo "After shutdown requests are ignored."
	@echo ""
	@echo "Stop existing server, if any -- command may error --"
	-@curl -i -X GET 'http://localhost:8123/api/shutdown' >/dev/null 2>/dev/null
	@rm -f out/test3*.out
	@echo "Compile and start new server"
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
	@sleep 7		# give it plenty of time to shutdown
	echo "Verify that the service is shutdown"
	./test5.sh
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
	@echo PASS

test4:
	echo "Verify correct value for pasword hash"
	@rm -f out/test4*.out
	./cleanup-this-one.sh
	curl -i -X GET 'http://localhost:8123/?password=angryMonkey' >out/test4.a.out 2>/dev/null 
	tail -1 out/test4.a.out >out/test4.b.out
	diff -w ref/test4.a.out out/test4.b.out
	@echo PASS

test5:
	echo "Test sending a SIGNAL to the service shut it down.  Currently uses HUP signal."
	@rm -f out/test4*.out out/pid-of-svr
	./cleanup-this-one.sh
	kill -HUP `cat out/pid-of-svr`
	echo "Verify that the service is shutdown"
	./test5.sh
	@echo PASS


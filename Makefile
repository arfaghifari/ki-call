build:
	@echo "building & running ki call"
	@go build -v -o backend_service main.go  && ./backend_service 

run:
	@echo "running ki call"
	@./backend_service 

gen-kitex:
	@echo "generate kitex"
	@kitex -module github.com/arfaghifari/ki-call proto/merchantvouchers.proto
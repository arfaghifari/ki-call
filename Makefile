build:
	@echo "building & running ki call"
	@go build -v -o backend_service main.go  && ./backend_service 

run:
	@echo "running ki call"
	@./backend_service 

gen-kitex:
	@echo "generate kitex $(obj)"
	@kitex -module github.com/arfaghifari/ki-call proto/$(obj).proto

gen-client:
	@echo "generate kitex client"
	@python3 generate_pkg.py
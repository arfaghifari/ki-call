build:
	@echo "building & running guild board"
	@go build -v -o backend_service main.go  && ./backend_service 

run:
	@echo "running guild board"
	@./backend_service 
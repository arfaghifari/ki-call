build:
	@echo "building & running guild board"
	@go build -v -o backend_service app.go  && ./backend_service 

run:
	@echo "running guild board"
	@./backend_service 
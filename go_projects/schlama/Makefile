install: build
	@go install github.com

install_win: build_win
	@go install github.com

run: build
	@./bin/schlama

run_win: build_win
	@./bin/schlama.exe

build: 
	@go build -o bin/schlama .

build_win:
	@go build -o bin/schlama.exe .

clean:
	@rm -Recurse ./bin

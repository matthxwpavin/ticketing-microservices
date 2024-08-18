dir := ./cmd/server
bin := ./server

# TODO: prequisite 'build' by $(bin) and $(bin) by all Go files of the project.
server:
	@echo "Building..."
	@go build -o $(bin) $(dir)
	@echo $(bin)

clean:
	@echo Cleaning up...
	@rm $(bin)
	@echo "$(bin) removed."

run: server
	@./env.sh
	@./server


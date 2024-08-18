dir := ./cmd/server
bin := ./server

# TODO: prequisite 'build' by $(bin) and $(bin) by all Go files of the project.
server:
	@echo "Building..."
	@go build -o $(bin) $(dir)
	@echo $(bin)

cleaning := Cleaning up...
removed := "$(bin) removed."
clean:
	@echo $(cleaning)
	@rm $(bin)
	@echo $(removed)

run: server
	@./run.sh


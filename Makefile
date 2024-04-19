default:
	@./scripts/build/build.sh
install:
	@echo "Installing.."
	@sudo cp bin/stellarpods /usr/local/bin/
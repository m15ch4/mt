alpha:
	./scripts/alpha.sh

alpha_run:
	go run micze.io/mt/cmd/alpha

clean:
	rm -rf ./bin

.PHONY: alpha alpha_run clean
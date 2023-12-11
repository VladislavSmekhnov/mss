run: build_jar
	docker compose up -d

build_jar:
	cd api-gateway && make rebuild

stop:
	docker compose down

deep-clean:
	docker system prune -a

.PHONY: build_jar run stop deep-clean

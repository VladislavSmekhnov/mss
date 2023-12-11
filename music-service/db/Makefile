IMAGE = img-db-mongo-music-service-lad:1.0
CONTAINER = db-mongo-music-service-lad

all: build_image run_container

build_image:
	docker build -t $(IMAGE) .

run_container:
	docker run -d -p 27017:27017 --name $(CONTAINER) $(IMAGE)

remove_image:
	docker rmi $(IMAGE)

stop_container:
	docker stop $(CONTAINER)

remove_container:
	docker rm $(CONTAINER)

rebuild: stop_container remove_container remove_image build_image run_container

clean:
	docker system prune -a

.PHONY: all build_image run_container remove_image stop_container remove_container rebuild clean

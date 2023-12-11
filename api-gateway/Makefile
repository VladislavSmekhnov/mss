JAR_NAME = api_gateway

build_jar:
	./mvnw clean package -DskipTests -Djar.finalName=$(JAR_NAME)

clean:
	rm -rf ./target/*.jar ./target/*.jar.original

rebuild: clean build_jar
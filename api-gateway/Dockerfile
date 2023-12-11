FROM openjdk:19-jdk-alpine

# Set the working directory in the container
WORKDIR /api-gateway-src

# Copy the compiled JAR file into the container at /app
COPY api-gateway/target/*.jar /api-gateway-src/api-gateway.jar

# Specify the command to run your application
CMD ["java", "-jar", "api-gateway.jar"]

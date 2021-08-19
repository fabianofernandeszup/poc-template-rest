FROM itau-ey5-docker.depdes.artifactory.prod.cloud.ihf/orangestack-baseimagejre:latest

EXPOSE 8080

RUN mkdir /app
ADD orangestack-*all.jar /app/$RIT_PARAMETER_APPLICATION_NAME.jar

ENTRYPOINT ["java", "-Djava.security.egd=file:/dev/./urandom", "-jar", "/app/$RIT_PARAMETER_APPLICATION_NAME.jar"]

pluginManagement {
    val artifactoryUrl: String by settings

    repositories {
        maven {
            url = uri("${artifactoryUrl}/gradle-devel")
        }
        maven {
            url = uri("${artifactoryUrl}/itau-ey5-tecnologia-gradle")
        }
    }
}
rootProject.name="$RIT_PARAMETER_APPLICATION_NAME"
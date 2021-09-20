val version : String by project
val group : String by project
val artifactoryUrl: String by project

plugins {
    id("com.itau.orangestack-plugin") version "0.1.0"
    id("io.micronaut.application") version "1.0.3"
}

dependencies {
    //Orangestack components
    implementation(platform("com.itau:orangestack-parent:0.3.1"))
    implementation("com.itau:orangestack-grpc")

    //tests
    testApi("com.itau:orangestack-tests")
    kaptTest("io.micronaut:micronaut-inject-java")
}

micronaut {
    runtime("netty")
    processing {
        val incremental = true
        val annotations = "${group}.*"
    }
}

application {
    mainClassName = "${group}.ApplicationKt"
}

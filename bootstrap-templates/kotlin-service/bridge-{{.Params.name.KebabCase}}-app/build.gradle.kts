dependencies {
    implementation("org.springframework.boot:spring-boot-starter")
    implementation("org.springframework.boot:spring-boot-starter-web")
    implementation("com.fasterxml.jackson.module:jackson-module-kotlin")
    implementation("org.springframework.boot:spring-boot-starter-data-jpa")

    runtimeOnly(kotlin("reflect"))

    implementation("com.bridgeapp:bridge-spring-boot-starter") {
        exclude(group = "org.springframework.boot", module = "spring-boot-starter-data-redis")
    }
    runtimeOnly("io.micrometer:micrometer-registry-statsd")
    runtimeOnly("org.slf4j:jcl-over-slf4j")
    runtimeOnly("org.springframework.boot:spring-boot-starter-actuator")

    runtimeOnly(project(":bridge-{{.Params.name}}-web"))

    testImplementation("org.springframework.boot:spring-boot-starter-test")
    testImplementation("org.springframework.boot:spring-boot-starter-data-jpa")
    testImplementation("com.bridgeapp:bridge-spring-boot-starter") {
        exclude(group = "org.springframework.boot", module = "spring-boot-starter-data-redis")
    }

    testImplementation("com.nhaarman:mockito-kotlin")
    testImplementation(project(":bridge-{{.Params.name}}-web"))

}

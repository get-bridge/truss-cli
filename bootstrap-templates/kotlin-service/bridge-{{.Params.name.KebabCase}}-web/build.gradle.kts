dependencies {
    implementation("org.springframework.boot:spring-boot-starter-web")
    implementation("org.springframework.data:spring-data-commons")
    implementation("io.springfox:springfox-swagger2")
    implementation("io.springfox:springfox-swagger-ui")
    // implementation("com.fasterxml.jackson.module:jackson-module-kotlin")

    runtimeOnly("org.slf4j:jcl-over-slf4j")

    testImplementation("org.springframework.boot:spring-boot-starter-test")
    testImplementation("com.nhaarman:mockito-kotlin")
}

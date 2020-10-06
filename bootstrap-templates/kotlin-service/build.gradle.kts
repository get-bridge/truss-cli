import org.springframework.boot.gradle.plugin.SpringBootPlugin

plugins {
    // Plugin versions are specified in settings.gradle.kts
    kotlin("jvm")
    kotlin("plugin.spring")
    kotlin("plugin.noarg")
    application
    id("com.gorylenko.gradle-git-properties")
    id("org.springframework.boot")
    id("io.spring.dependency-management")
    id("co.uzzu.dotenv.gradle")
}

subprojects {
    apply(plugin = "java-library")
    apply(plugin = "kotlin")
    apply(plugin = "kotlin-spring")
    apply(plugin = "kotlin-jpa")
    apply(plugin = "io.spring.dependency-management")
}

allprojects {
    group = "com.bridgeapp.{{.Params.name.FlatCase}}"
    version = "0.0.1-SNAPSHOT"

    java {
        sourceCompatibility = JavaVersion.VERSION_11 // If targetCompatibility is not set it defaults to sourceCompatibility
    }

    kotlin {
        target.compilations.all {
            kotlinOptions {
                javaParameters = true
                jvmTarget = project.java.targetCompatibility.majorVersion
            }
        }
    }

    repositories {
        mavenCentral()

        maven {
            name = "Instructure Artifactory"
            url = uri("https://instructure.jfrog.io/artifactory/internal-maven")

            credentials {
                val artifactoryUsername: String by project
                val artifactoryPassword: String by project

                username = artifactoryUsername
                password = artifactoryPassword
            }
        }
    }

    dependencyManagement {
        dependencies {
            // Dependency versions are defined in gradle.properties
            val bridgeSpringBootVersion: String by project
            val mockitoKotlinVersion: String by project
            val springFoxSwaggerVersion: String by project

            imports {
                mavenBom(SpringBootPlugin.BOM_COORDINATES)
                mavenBom("com.bridgeapp:bridge-spring-boot-starter:$bridgeSpringBootVersion")
            }

            dependencySet("com.bridgeapp:$bridgeSpringBootVersion") {
                entry("bridge-spring-boot-starter")
                entry("bridge-messagebus")
            }
            dependency("com.nhaarman:mockito-kotlin:$mockitoKotlinVersion")
            dependencySet("io.springfox:$springFoxSwaggerVersion") {
                entry("springfox-swagger2")
                entry("springfox-swagger-ui")
            }
        }
    }

    configurations {
        all {
            // Unneeded logging dependencies, usually from spring-boot-starter
            exclude(group = "org.apache.logging.log4j")
            exclude(group = "org.springframework", module = "spring-jcl")

            // Unneeded test dependencies, usually from spring-boot-starter-test
            exclude(group = "com.jayway.jsonpath", module = "json-path")
            exclude(group = "org.assertj", module = "assertj-core")
            exclude(group = "org.hamcrest", module = "hamcrest")
            exclude(group = "org.junit.vintage")
            exclude(group = "org.skyscreamer", module = "jsonassert")
            exclude(group = "org.xmlunit", module = "xmlunit-core")
        }
    }
}

dependencies {
    runtimeOnly(project(":bridge-{{.Params.name}}-app"))
}

application {
    mainClassName = "com.bridgeapp.{{.Params.name.FlatCase}}.{{.Params.name.PascalCase}}Main"
}

tasks {
    bootRun {
        systemProperty("spring.profiles.include", "dev,local")
        systemProperty("spring.config.additional-location", "file://${rootProject.file("config")}/")

        environment(env.allVariables)
    }

    bootJar {
        isPreserveFileTimestamps = false
        isReproducibleFileOrder = true
    }
}

subprojects {
    dependencies {
        implementation(kotlin("stdlib-jdk8"))
    }

    tasks {
        jar {
            isPreserveFileTimestamps = false
            isReproducibleFileOrder = true
        }

        test {
            useJUnitPlatform()

            systemProperty("spring.profiles.include", "test,test-local")
            systemProperty("spring.config.additional-location", "file://${rootProject.file("config")}/")

            environment(env.allVariables)
        }
    }
}

val printPropertyTask = tasks.create("printProperty") {
    project.findProperty("propertyName")?.let { propertyName ->
        project.findProperty(propertyName.toString())?.let { propertyValue ->
            logger.quiet("{}", propertyValue)
        }
    }
}

printPropertyTask.apply {
    group = "Help"
    description = "Print value for a given property (-PpropertyName=PROPERTY_NAME)"
}

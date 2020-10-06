pluginManagement {
    // Versions are defined in gradle.properties
    val dependencyManagementPluginVersion: String by settings
    val dotEnvPluginVersion: String by settings
    val gitPropertiesPluginVersion: String by settings
    val kotlinVersion = settings.extra["kotlin.version"]?.toString()
    val springBootVersion: String by settings

    plugins {
        kotlin("jvm") version kotlinVersion
        kotlin("plugin.noarg") version kotlinVersion
        kotlin("plugin.spring") version kotlinVersion
        id("co.uzzu.dotenv.gradle") version dotEnvPluginVersion
        id("com.gorylenko.gradle-git-properties") version gitPropertiesPluginVersion
        id("io.spring.dependency-management") version dependencyManagementPluginVersion
        id("org.springframework.boot") version springBootVersion
    }
}

rootProject.name = "bridge-{{.Params.name}}-backend"

include(":bridge-{{.Params.name}}-app")
include(":bridge-{{.Params.name}}-web")
include(":bridge-{{.Params.name}}-web-test")

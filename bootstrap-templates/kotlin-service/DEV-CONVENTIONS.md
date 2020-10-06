# Development conventions

## Kotlin coding conventions

Follow general Kotlin coding conventions:

* https://kotlinlang.org/docs/reference/coding-conventions.html
* https://kotlinlang.org/docs/reference/idioms.html

## Dependency management

* Specify artifact versions in `gradle.properties` with lower camel case name: `artifactNameVersion` (e.g. `kotlinVersion`)
* Keep version properties in `gradle.properties` sorted by name within a block
* Define group-wide version instead of artifact-wide (`jacksonKotlinModuleVersion` (bad) vs. `jacksonVersion` (good))
* Specify dependency version in `dependencyManagement` block in top-level `build.gradle.kts`
* Keep dependencies sorted by id (`group:name`) within configurations
* Define `project()` dependencies at the last of the dependency list within a configuration
* Specify dependencies as `group:name` string instead of explicit group = "", name = "" 
* Use Spring Boot BOM-provided version of a dependency if that's defined in there
* Import BOM of the newly required dependency if there is one instead of the concrete artifact
* If a different version is needed, override the version property coming from Spring Boot BOM (e.g. `junit-jupiter.version` in `gradle.properties`)
* Prefer `runtimeOnly` (and `testRuntimeOnly`) configuration over `implementation` (and `testImplementation`)
* Prefer `implementation` configuration over `api` in library projects (currently all subproject)

## Spring application properties

* Prefer YAML over properties (`application.yml` over `application.properties`)
* Use __kebab-case__ for property names

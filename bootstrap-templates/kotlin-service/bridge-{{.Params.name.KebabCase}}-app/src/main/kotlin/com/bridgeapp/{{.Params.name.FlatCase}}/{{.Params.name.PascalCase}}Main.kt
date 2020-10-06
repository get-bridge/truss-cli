@file:JvmName("{{.Params.name.PascalCase}}Main")

package com.bridgeapp.{{.Params.name.FlatCase}}

import com.bridgeapp.autoconfigure.multitenant.MultiTenantAutoConfiguration
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.autoconfigure.domain.EntityScan
import org.springframework.boot.runApplication
import org.springframework.data.jpa.repository.config.EnableJpaRepositories

@SpringBootApplication(
    scanBasePackages=["com.bridgeapp.{{.Params.name.FlatCase}}"],
    exclude = [
        MultiTenantAutoConfiguration::class
    ])
@EntityScan(basePackages=["com.bridgeapp.{{.Params.name.FlatCase}}"])
@EnableJpaRepositories(basePackages=["com.bridgeapp.{{.Params.name.FlatCase}}"])
class {{.Params.name.PascalCase}}App

fun main(vararg args: String) {
    runApplication<{{.Params.name.PascalCase}}App>(*args)
}

@file:JvmName("{{.Params.name.PascalCase}}Main")

package com.bridgeapp.{{.Params.name.FlatCase}}

import com.bridgeapp.autoconfigure.multitenant.HibernateMultiTenantAutoConfiguration
import com.bridgeapp.autoconfigure.multitenant.MultiTenantAutoConfiguration
import com.bridgeapp.autoconfigure.multitenant.TenantServiceAutoConfiguration
import org.springframework.boot.autoconfigure.SpringBootApplication

@SpringBootApplication(exclude = [
    MultiTenantAutoConfiguration::class,
    TenantServiceAutoConfiguration::class,
    HibernateMultiTenantAutoConfiguration::class
])
class {{.Params.name.PascalCase}}TestApp

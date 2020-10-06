package com.bridgeapp.{{.Params.name.FlatCase}}

import com.bridgeapp.autoconfigure.multitenant.MultiTenantProperties
import com.bridgeapp.autoconfigure.multitenant.TenantServiceAutoConfiguration
import com.bridgeapp.multitenant.TenantService
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.autoconfigure.AutoConfigureAfter
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.context.annotation.Profile
import org.springframework.web.servlet.handler.MappedInterceptor

/**
 * Configuration for Multi-tenancy
 */
@Profile("!test")
@Configuration
@AutoConfigureAfter(
    TenantServiceAutoConfiguration::class
)
@EnableConfigurationProperties(MultiTenantProperties::class)
class FlexibleMultiTenantAutoConfiguration @Autowired constructor(
    private val multiTenantProperties: MultiTenantProperties
) {
    /**
     * Initialize HTTP interceptor to parse tenant for requests
     * Ignores tenant-less requests (actuator, etc)
     */
    @Bean
    fun mappedInterceptor(
        tenantService: TenantService
    ) =
        MappedInterceptor(arrayOf("/**"),
            multiTenantProperties.allExcludePatterns,
            FlexibleTenantInterceptor(multiTenantProperties.hostName, tenantService))
}

package com.bridgeapp.{{.Params.name.FlatCase}}

import com.bridgeapp.multitenant.TenantResolver
import com.bridgeapp.multitenant.TenantService
import javax.servlet.http.HttpServletRequest
import javax.servlet.http.HttpServletResponse
import org.springframework.web.servlet.handler.HandlerInterceptorAdapter

/**
 * Intercept HTTP requests then parse the current tenant
 *
 * This version creates the tenant on the fly if it does not exists
 *
 * Checks for x-tenant in the request header, if not found
 * checks the subdomain name
 *
 * Must be given to spring mapped interceptor:
 * <pre>
 *   {@code
 *       @Bean
 *       fun mappedInterceptor(schemaRepository: AbstractSchemaRepository) =
 *         MappedInterceptor(
 *           ...
 *           TenantInterceptor(hostName, schemaRepository)
 *         )
 *   }
 * </pre>
 */
class FlexibleTenantInterceptor(
    private val hostName: String,
    private val tenantService: TenantService
) : HandlerInterceptorAdapter() {

    @Throws(Exception::class)
    override fun preHandle(
        request: HttpServletRequest,
        response: HttpServletResponse,
        handler: Any
    ): Boolean {
        val tenantId = parseTenantId(request)
        if (!tenantService.exists(tenantId)) {
            tenantService.create(tenantId)
        }
        TenantResolver.tenantId = tenantId
        return true
    }

    private fun parseTenantId(request: HttpServletRequest): String {
        var tenantId = request.getHeader("x-tenant")
        if (tenantId.isNullOrEmpty()) {
            tenantId = request.serverName.split(".$hostName").first()
        }
        return tenantId
    }
}

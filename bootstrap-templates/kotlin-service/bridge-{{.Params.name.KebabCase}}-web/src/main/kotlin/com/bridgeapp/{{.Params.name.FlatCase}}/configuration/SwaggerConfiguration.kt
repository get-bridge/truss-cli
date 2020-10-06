package com.bridgeapp.{{.Params.name.FlatCase}}.configuration

import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import springfox.documentation.builders.ApiInfoBuilder
import springfox.documentation.builders.PathSelectors
import springfox.documentation.builders.RequestHandlerSelectors
import springfox.documentation.service.Contact
import springfox.documentation.spi.DocumentationType
import springfox.documentation.spring.web.plugins.Docket
import springfox.documentation.swagger2.annotations.EnableSwagger2

@Configuration
@EnableSwagger2
class SwaggerConfiguration {
    private val apiInfo = ApiInfoBuilder()
        .title("Bridge {{.Params.name}} API")
        .description("Send questions to #voyager on Slack")
        .version("1.0.0")
        .contact(Contact("Voyager", null, null))
        .build()

    @Bean
    fun docket(): Docket {
        return Docket(DocumentationType.SWAGGER_2)
            .apiInfo(apiInfo)
            .select()
            .apis(RequestHandlerSelectors.basePackage("com.bridgeapp"))
            .paths(PathSelectors.any())
            .build()
    }
}

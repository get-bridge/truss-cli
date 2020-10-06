package com.bridgeapp.{{.Params.name.FlatCase}}

import io.swagger.annotations.ApiOperation
import io.swagger.annotations.ApiResponse
import io.swagger.annotations.ApiResponses
import org.springframework.http.MediaType
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController

@RestController
@RequestMapping
class {{.Params.name.PascalCase}}Controller() {
    @ApiOperation("Index Response Example")
    @GetMapping
    fun index(): ResponseEntity<String> = ResponseEntity.ok("Hello, World!")

    @ApiOperation("Api Response Example")
    @GetMapping("/hello", produces = [MediaType.APPLICATION_JSON_VALUE])
    fun hello(): ResponseEntity<Hello> = ResponseEntity.ok(Hello("Hello, World!"))

    @ApiOperation("Parameterized Api Response Example")
    @PostMapping("/hello", produces = [MediaType.APPLICATION_JSON_VALUE])
    fun hello(person: Person): ResponseEntity<Hello> = ResponseEntity.ok(Hello("Hello, ${person.name}!"))
}

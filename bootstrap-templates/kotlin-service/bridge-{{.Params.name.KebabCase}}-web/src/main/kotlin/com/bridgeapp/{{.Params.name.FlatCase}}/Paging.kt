package com.bridgeapp.{{.Params.name.FlatCase}}

import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Pageable
import org.springframework.data.domain.Sort

private const val DEFAULT_PAGE_NUM = 0
private const val DEFAULT_PAGE_SIZE = 20

internal fun getPageRequest(page: Int?, pageSize: Int?, sort: String?): Pageable =
    if (sort != null) {
        PageRequest.of(page ?: DEFAULT_PAGE_NUM, pageSize ?: DEFAULT_PAGE_SIZE, Sort.by(sort))
    } else {
        PageRequest.of(page ?: DEFAULT_PAGE_NUM, pageSize ?: DEFAULT_PAGE_SIZE)
    }

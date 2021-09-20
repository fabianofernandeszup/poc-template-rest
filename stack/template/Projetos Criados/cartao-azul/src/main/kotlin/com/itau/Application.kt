package com.itau

import io.micronaut.runtime.Micronaut.*
fun main(args: Array<String>) {
	build()
	    .args(*args)
		.packages("$RIT_PARAMETER_APPLICATION_PACKAGE")
		.start()
}
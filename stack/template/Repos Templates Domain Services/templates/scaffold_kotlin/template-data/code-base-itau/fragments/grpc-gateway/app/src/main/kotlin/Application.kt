import io.micronaut.runtime.Micronaut.build

// package $RIT_PARAMETER_APPLICATION_PACKAGE

fun main(args: Array<String>) {
    build()
        .args(*args)
        .packages("$RIT_PARAMETER_APPLICATION_PACKAGE")
        .start()
}
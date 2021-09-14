package com.itau.adapters.kafka.consumer

import com.itau.avro.Product
import com.itau.orangestack.logging.LoggableClass
import io.micronaut.configuration.kafka.annotation.*
import io.micronaut.messaging.Acknowledgement

@KafkaListener(
    groupId = "group-pubsub-consumer-group",
    offsetReset = OffsetReset.EARLIEST,
    offsetStrategy = OffsetStrategy.DISABLED
)
class ProductListener {
    @Topic("kaas-pubsub-demo")
    fun receive(@KafkaKey brand: String, product: Product?, acknowledgement: Acknowledgement) {
        logger.info { "Got Product - ${product?.name} by $brand" }
        acknowledgement.ack()
    }

    companion object : LoggableClass()
}
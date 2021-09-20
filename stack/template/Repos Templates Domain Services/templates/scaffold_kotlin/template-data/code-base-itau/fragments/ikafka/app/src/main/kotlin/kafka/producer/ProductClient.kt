package com.itau.adapters.kafka.producer

import com.itau.avro.Product
import io.micronaut.configuration.kafka.annotation.KafkaClient
import io.micronaut.configuration.kafka.annotation.KafkaKey
import io.micronaut.configuration.kafka.annotation.Topic

@KafkaClient(id = "product-client")
interface ProductClient {
    @Topic("kaas-pubsub-demo")
    fun sendProduct(@KafkaKey brand: String?, product: Product)
    fun sendProduct(@Topic topic: String?, @KafkaKey brand: String?, product: Product)
}
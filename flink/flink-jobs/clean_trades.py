from pyflink.common import Types
from pyflink.datastream import StreamExecutionEnvironment, OutputTag
from pyflink.datastream.connectors import FlinkKafkaConsumer, FlinkKafkaProducer
from pyflink.datastream.functions import ProcessFunction, RuntimeContext
from pyflink.datastream.formats.json import JsonRowDeserializationSchema, JsonRowSerializationSchema

env = StreamExecutionEnvironment.get_execution_environment()
bootstrap_servers = "redpanda-1:29092,redpanda-2:29093,redpanda-3:29094"

# JSON schema definition for deserialization
json_deserialization_schema = JsonRowDeserializationSchema.builder() \
    .type_info(
        type_info=Types.ROW_NAMED(
            ["base", "quote", "exchange", "volume_base", "volume_quote", "price", "timestamp"],
            [Types.STRING(), Types.STRING(), Types.STRING(), Types.DOUBLE(), Types.DOUBLE(), Types.DOUBLE(), Types.STRING()]
        )
    ).build()

# JSON schema definition for serialization
json_serialization_schema = JsonRowSerializationSchema.builder() \
    .with_type_info(
        type_info=Types.ROW_NAMED(
            ["base", "quote", "exchange", "volume_base", "volume_quote", "price", "timestamp"],
            [Types.STRING(), Types.STRING(), Types.STRING(), Types.DOUBLE(), Types.DOUBLE(), Types.DOUBLE(), Types.STRING()]
        )
    ).build()

# Source: Kafka Consumer
kafka_consumer = FlinkKafkaConsumer(
    topics="coincap-crypto",
    deserialization_schema=json_deserialization_schema,
    properties={"bootstrap.servers": bootstrap_servers, "group.id": "trade-group"},
)

# Output tags for valid and invalid messages
valid_msg_output_tag = OutputTag(
    "valid-trades",
    Types.ROW_NAMED(
        ["base", "quote", "exchange", "volume_base", "volume_quote", "price", "timestamp"],
        [Types.STRING(), Types.STRING(), Types.STRING(), Types.DOUBLE(), Types.DOUBLE(), Types.DOUBLE(), Types.STRING()]
    )
)
invalid_msg_output_tag = OutputTag(
    "invalid-trades",
    Types.ROW_NAMED(
        ["base", "quote", "exchange", "volume_base", "volume_quote", "price", "timestamp"],
        [Types.STRING(), Types.STRING(), Types.STRING(), Types.DOUBLE(), Types.DOUBLE(), Types.DOUBLE(), Types.STRING()]
    )
)

class ValidateTrade(ProcessFunction):

    def process_element(self, value, ctx: RuntimeContext):
        # Assuming 'value' is a Row with the fields you mentioned
        valid = True
        for field in value:
            if isinstance(field, str) and (field is None or field == ""):
                valid = False
            elif isinstance(field, (float, int)) and field <= 0:
                valid = False
        
        if valid:
            yield valid_msg_output_tag, value
        else:
            yield invalid_msg_output_tag, value

# Apply the process function to filter and split the stream
processed_stream = env.add_source(kafka_consumer).process(ValidateTrade())

# Sink: Kafka Producer for messages with valid fields
valid_trades_producer = FlinkKafkaProducer(
    topic="valid-trades",
    serialization_schema=json_serialization_schema,
    producer_config={"bootstrap.servers": bootstrap_servers}
)

# Sink: Kafka Producer for messages with invalid fields
invalid_trades_producer = FlinkKafkaProducer(
    topic="invalid-trades",
    serialization_schema=json_serialization_schema,
    producer_config={"bootstrap.servers": bootstrap_servers}
)

# Write valid and invalid messages to respective Kafka topics
processed_stream.get_side_output(valid_msg_output_tag).add_sink(valid_trades_producer)
processed_stream.get_side_output(invalid_msg_output_tag).add_sink(invalid_trades_producer)

env.execute("Trade Validation Job")



# -------------------------------------------------
# trade_stream = env.add_source(kafka_consumer)
# trade_stream.print()
# env.execute("Print Trades Stream")

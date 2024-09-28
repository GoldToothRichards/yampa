CREATE DATABASE IF NOT EXISTS trades;

-- LowCardinality should be used for columns with fewer than ~10k-100k distinct values

CREATE TABLE IF NOT EXISTS trades.raw_trades(
    source LowCardinality(String) COMMENT 'Source of the data, i.e. websocket origin.' CODEC(ZSTD),
    base LowCardinality(String) COMMENT 'The base asset, i.e. Bitcoin.' CODEC(ZSTD),
    quote LowCardinality(String) COMMENT 'The quote asset, i.e. USD.' CODEC(ZSTD),
    exchange LowCardinality(String) COMMENT 'The exchange where the trade occured, i.e. Coinbase.' CODEC(ZSTD),
    volume_base Float64 COMMENT 'The amount of the base asset exchanged in the trade, measured in base units.' CODEC(ZSTD),
    volume_quote Float64 COMMENT 'The amount of the quote asset exchanged in the trade, measured in quote units.' CODEC(ZSTD),
    price Float64 COMMENT 'The price of the base asset in terms of the quote asset for the given trade, i.e. the ratio (volume_quote/volume_base).' CODEC(ZSTD),
    timestamp DateTime64(3) COMMENT 'The timestamp of the trade, as provided by the source. UTC timezone, millisecond precision.' CODEC(ZSTD),
)
ENGINE = MergeTree
PRIMARY KEY (base, quote, exchange, source)
ORDER BY (base, quote, exchange, source, timestamp)
PARTITION BY toYYYYMM(timestamp)
TTL
    toDateTime(timestamp) + INTERVAL 1 MONTH TO VOLUME 's3',
    toDateTime(timestamp) + INTERVAL 1 YEAR DELETE
SETTINGS storage_policy = 'shared';

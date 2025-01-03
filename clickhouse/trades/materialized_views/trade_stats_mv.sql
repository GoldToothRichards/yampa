CREATE DATABASE IF NOT EXISTS trades;

CREATE TABLE IF NOT EXISTS trades.stats_24h(
    base LowCardinality(String) COMMENT 'The base asset, i.e. Bitcoin.' CODEC(ZSTD),
    quote LowCardinality(String) COMMENT 'The quote asset, i.e. USD.' CODEC(ZSTD),
    window_start String COMMENT 'The start timestamp of the stat collection window.' CODEC(ZSTD),
    window_end String COMMENT 'The end timestamp of the stat collection window' CODEC(ZSTD),
    avg_volume_base Float64 COMMENT 'The average amount of the base asset exchanged in the trades within the window' CODEC(ZSTD),
    avg_volume_quote Float64 COMMENT 'The average amount of the quote asset exchanged in the trades within the window' CODEC(ZSTD),
    avg_price Float64 COMMENT 'The (simple) average price of the asset for trades within the window.' CODEC(ZSTD),
    stddev_volume_base Float64 COMMENT 'The standard deviation of the volume of the base asset for trades within the window' CODEC(ZSTD),
    stddev_volume_quote Float64 COMMENT 'The standard deviation of the volume of the quote asset for trades within the window' CODEC(ZSTD),
    stddev_price Float64 COMMENT 'The standard deviation of the price of the asset for trades within the window.' CODEC(ZSTD)
)
ENGINE = MergeTree()
PRIMARY KEY (base, quote)
ORDER BY (base, quote, window_start);

CREATE MATERIALIZED VIEW IF NOT EXISTS trades.stats_mv_24h TO trades.stats_24h AS
    SELECT
        base,
        quote,
        tumbleStart(parseDateTimeBestEffort(timestamp), INTERVAL 1 DAY) AS window_start,
        tumbleEnd(parseDateTimeBestEffort(timestamp), INTERVAL 1 DAY) AS window_end,
        avg(volume_base) AS avg_volume_base,
        avg(volume_quote) AS avg_volume_quote,
        avg(price) AS avg_price,
        stddevPop(volume_base) AS stddev_volume_base,
        stddevPop(volume_quote) AS stddev_volume_quote,
        stddevPop(price) AS stddev_price
    FROM trades.raw_trades
    GROUP BY base, quote, window_start, window_end;


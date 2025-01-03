CREATE DATABASE IF NOT EXISTS trades;

CREATE TABLE IF NOT EXISTS trades.stats_24h(
    base LowCardinality(String) COMMENT 'The base asset, i.e. Bitcoin.' CODEC(ZSTD),
    quote LowCardinality(String) COMMENT 'The quote asset, i.e. USD.' CODEC(ZSTD),
    window_start DateTime COMMENT 'The start timestamp of the stat collection window.' CODEC(ZSTD),
    window_end DateTime COMMENT 'The end timestamp of the stat collection window' CODEC(ZSTD),
    avg_volume_base AggregateFunction(avg, Float64) COMMENT 'The average amount of the base asset exchanged in the trades within the window' CODEC(ZSTD),
    avg_volume_quote AggregateFunction(avg, Float64) COMMENT 'The average amount of the quote asset exchanged in the trades within the window' CODEC(ZSTD),
    avg_price AggregateFunction(avg, Float64) COMMENT 'The (simple) average price of the asset for trades within the window.' CODEC(ZSTD),
    stddev_volume_base AggregateFunction(stddevPop, Float64) COMMENT 'The standard deviation of the volume of the base asset for trades within the window' CODEC(ZSTD),
    stddev_volume_quote AggregateFunction(stddevPop, Float64) COMMENT 'The standard deviation of the volume of the quote asset for trades within the window' CODEC(ZSTD),
    stddev_price AggregateFunction(stddevPop, Float64) COMMENT 'The standard deviation of the price of the asset for trades within the window.' CODEC(ZSTD)
)
ENGINE = AggregatingMergeTree()
PRIMARY KEY (base, quote)
ORDER BY (base, quote, window_start);

CREATE MATERIALIZED VIEW IF NOT EXISTS trades.stats_mv_24h TO trades.stats_24h AS
    SELECT
        base,
        quote,
        hopStart(timestamp, INTERVAL 10 MINUTE, INTERVAL 1 DAY) AS window_start,
        hopEnd(timestamp, INTERVAL 10 MINUTE, INTERVAL 1 DAY) AS window_end,
        avgState(volume_base) AS avg_volume_base,
        avgState(volume_quote) AS avg_volume_quote,
        avgState(price) AS avg_price,
        stddevPopState(volume_base) AS stddev_volume_base,
        stddevPopState(volume_quote) AS stddev_volume_quote,
        stddevPopState(price) AS stddev_price
    FROM trades.raw_trades
    GROUP BY base, quote, window_start, window_end;


DROP TABLE trades.stats_24h;
DROP VIEW trades.stats_mv_24h;

SELECT COUNT(*) FROM trades.raw_trades;
SELECT
    base,
    quote,
    window_start,
    window_end,
    avgMerge(avg_volume_base) AS avg_volume_base,
    avgMerge(avg_volume_quote) AS avg_volume_quote,
    avgMerge(avg_price) AS avg_price,
    stddevPopMerge(stddev_volume_base) AS stddev_volume_base,
    stddevPopMerge(stddev_volume_quote) AS stddev_volume_quote,
    stddevPopMerge(stddev_price) AS stddev_price
FROM trades.stats_24h
GROUP BY
    base,
    quote,
    window_start,
    window_end
ORDER BY
    base,
    quote
LIMIT 5
-- Table for 1m price ticks
CREATE TABLE IF NOT EXISTS trades.price_ticks_1m (
    base LowCardinality(String),
    quote LowCardinality(String),
    exchange LowCardinality(String),
    timestamp DateTime,
    open Float64,
    high Float64,
    low Float64,
    close Float64,
    volume_base Float64,
    volume_quote Float64
)
ENGINE = AggregatingMergeTree()
ORDER BY (base, quote, exchange, timestamp)
PARTITION BY toYYYYMM(timestamp);

-- 1m price ticks
CREATE MATERIALIZED VIEW trades.price_ticks_1m_mv TO trades.price_ticks_1m AS 
WITH clean_trades AS
(
    SELECT * FROM trades.raw_trades
    WHERE source != ''
      AND base != ''
      AND quote != ''
      AND exchange != ''
      AND price > 0
      AND volume_base > 0
      AND volume_quote > 0
      AND base != quote
)
SELECT
    base,
    quote,
    exchange,
    toStartOfInterval(timestamp, INTERVAL 1 MINUTE) AS timestamp,
    argMin(price, timestamp) AS open,
    max(price) AS high,
    min(price) AS low,
    argMax(price, timestamp) AS close,
    sum(volume_base) AS volume_base,
    sum(volume_quote) AS volume_quote
FROM clean_trades
GROUP BY
    base,
    quote,
    exchange,
    timestamp;

-- Table for 5m price ticks
CREATE TABLE IF NOT EXISTS trades.price_ticks_5m (
    base LowCardinality(String),
    quote LowCardinality(String),
    exchange LowCardinality(String),
    timestamp DateTime,
    open Float64,
    high Float64,
    low Float64,
    close Float64,
    volume_base Float64,
    volume_quote Float64
)
ENGINE = AggregatingMergeTree()
ORDER BY (base, quote, exchange, timestamp);

-- 5m price ticks
CREATE MATERIALIZED VIEW trades.price_ticks_5m_mv TO trades.price_ticks_5m
POPULATE AS SELECT
    base,
    quote,
    exchange,
    toStartOfInterval(timestamp, INTERVAL 5 MINUTE) AS timestamp,
    argMin(open, timestamp) AS open,
    max(high) AS high,
    min(low) AS low,
    argMax(close, timestamp) AS close,
    sum(volume_base) AS volume_base,
    sum(volume_quote) AS volume_quote
FROM trades.price_ticks_1m_mv
GROUP BY
    base,
    quote,
    exchange,
    timestamp;

-- Table for 15m price ticks
CREATE TABLE IF NOT EXISTS trades.price_ticks_15m (
    base LowCardinality(String),
    quote LowCardinality(String),
    exchange LowCardinality(String),
    timestamp DateTime,
    open Float64,
    high Float64,
    low Float64,
    close Float64,
    volume_base Float64,
    volume_quote Float64
)
ENGINE = AggregatingMergeTree()
ORDER BY (base, quote, exchange, timestamp);

-- 15m price ticks
CREATE MATERIALIZED VIEW trades.price_ticks_15m_mv TO trades.price_ticks_15m
POPULATE AS SELECT
    base,
    quote,
    exchange,
    toStartOfInterval(timestamp, INTERVAL 15 MINUTE) AS timestamp,
    argMin(open, timestamp) AS open,
    max(high) AS high,
    min(low) AS low,
    argMax(close, timestamp) AS close,
    sum(volume_base) AS volume_base,
    sum(volume_quote) AS volume_quote
FROM trades.price_ticks_5m_mv
GROUP BY
    base, quote, exchange, timestamp;

-- Table for 1h price ticks
CREATE TABLE IF NOT EXISTS trades.price_ticks_1h (
    base LowCardinality(String),
    quote LowCardinality(String),
    exchange LowCardinality(String),
    timestamp DateTime,
    open Float64,
    high Float64,
    low Float64,
    close Float64,
    volume_base Float64,
    volume_quote Float64
)
ENGINE = AggregatingMergeTree()
ORDER BY (base, quote, exchange, timestamp);

-- 1h price ticks
CREATE MATERIALIZED VIEW trades.price_ticks_1h_mv TO trades.price_ticks_1h
POPULATE AS SELECT
    base,
    quote,
    exchange,
    toStartOfInterval(timestamp, INTERVAL 1 HOUR) AS timestamp,
    argMin(open, timestamp) AS open,
    max(high) AS high,
    min(low) AS low,
    argMax(close, timestamp) AS close,
    sum(volume_base) AS volume_base,
    sum(volume_quote) AS volume_quote
FROM trades.price_ticks_5m_mv
GROUP BY
    base,
    quote,
    exchange,
    timestamp;

-- Table for 4h price ticks
CREATE TABLE IF NOT EXISTS trades.price_ticks_4h (
    base LowCardinality(String),
    quote LowCardinality(String),
    exchange LowCardinality(String),
    timestamp DateTime,
    open Float64,
    high Float64,
    low Float64,
    close Float64,
    volume_base Float64,
    volume_quote Float64
)
ENGINE = AggregatingMergeTree()
ORDER BY (base, quote, exchange, timestamp);

-- 4h price ticks
CREATE MATERIALIZED VIEW trades.price_ticks_4h_mv TO trades.price_ticks_4h
POPULATE AS SELECT
    base,
    quote,
    exchange,
    toStartOfInterval(timestamp, INTERVAL 4 HOUR) AS timestamp,
    argMin(open, timestamp) AS open,
    max(high) AS high,
    min(low) AS low,
    argMax(close, timestamp) AS close,
    sum(volume_base) AS volume_base,
    sum(volume_quote) AS volume_quote
FROM trades.price_ticks_1h_mv
GROUP BY
    base,
    quote,
    exchange,
    timestamp;

-- Table for 1d price ticks
CREATE TABLE IF NOT EXISTS trades.price_ticks_1d (
    base LowCardinality(String),
    quote LowCardinality(String),
    exchange LowCardinality(String),
    timestamp DateTime,
    open Float64,
    high Float64,
    low Float64,
    close Float64,
    volume_base Float64,
    volume_quote Float64
)
ENGINE = AggregatingMergeTree()
ORDER BY (base, quote, exchange, timestamp);

-- 1d price ticks
CREATE MATERIALIZED VIEW trades.price_ticks_1d_mv TO trades.price_ticks_1d
POPULATE AS SELECT
    base,
    quote,
    exchange,
    toStartOfInterval(timestamp, INTERVAL 1 DAY) AS timestamp,
    argMin(open, timestamp) AS open,
    max(high) AS high,
    min(low) AS low,
    argMax(close, timestamp) AS close,
    sum(volume_base) AS volume_base,
    sum(volume_quote) AS volume_quote
FROM trades.price_ticks_4h_mv
GROUP BY
    base,
    quote,
    exchange,
    timestamp;
-- Table for 1m price ticks
CREATE TABLE IF NOT EXISTS trades.price_ticks_1m (
    base LowCardinality(String),
    quote LowCardinality(String),
    exchange LowCardinality(String),
    source LowCardinality(String),
    window_start DateTime,
    window_end DateTime,
    open_state AggregateFunction(argMin, Float64, DateTime),
    high_state AggregateFunction(max, Float64),
    low_state AggregateFunction(min, Float64),
    close_state AggregateFunction(argMax, Float64, DateTime),
    volume_base_state AggregateFunction(sum, Float64),
    volume_quote_state AggregateFunction(sum, Float64)
)
ENGINE = AggregatingMergeTree()
PRIMARY KEY (base, quote, exchange, source)
ORDER BY (base, quote, exchange, source, window_start)
PARTITION BY toYYYYMM(window_start)
TTL
    window_start + INTERVAL 1 MONTH TO VOLUME 's3',
    window_start + INTERVAL 1 YEAR DELETE;

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
    source,
    tumbleStart(toDateTime(timestamp), INTERVAL 1 MINUTE) AS window_start,
    tumbleEnd(toDateTime(timestamp), INTERVAL 1 MINUTE) AS window_end,
    argMinState(price, toDateTime(timestamp)) AS open_state,
    maxState(price) AS high_state,
    minState(price) AS low_state,
    argMaxState(price, toDateTime(timestamp)) AS close_state,
    sumState(volume_base) AS volume_base_state,
    sumState(volume_quote) AS volume_quote_state
FROM clean_trades
GROUP BY
    base,
    quote,
    exchange,
    source,
    window_start,
    window_end;

-- Test params
SET param_start = '2024-09-28 18:00:00';
SET param_end = '2024-09-28 19:00:00';
SET param_base = 'bitcoin';
SET param_quote = 'tether';
SET param_exchange = 'gdax';

-- OHLCV View
CREATE OR REPLACE VIEW trades.ohlcv_1m AS
SELECT
    base,
    quote,
    exchange,
    source,
    window_start,
    argMinMerge(open_state) as open,
    maxMerge(high_state) as high,
    minMerge(low_state) as low,
    argMaxMerge(close_state) as close,
    sumMerge(volume_base_state) as volume_base,
    sumMerge(volume_quote_state) as volume_quote,
    (close - open) / open as pct_change,
    volume_quote / volume_base as vwap
FROM trades.price_ticks_1m
WHERE (base = {base: String}) AND (quote = {quote: String}) AND (exchange = {exchange: String})
AND (window_start >= {start: String}) AND (window_start <= {end: String})
GROUP BY
    base,
    quote,
    exchange,
    source,
    window_start
ORDER BY
    base,
    quote,
    exchange,
    source,
    window_start ASC
WITH FILL
    FROM toDateTime({start: String})
    TO toDateTime({end: String})
    STEP INTERVAL 1 MINUTE
    INTERPOLATE (vwap)

-- Test
SELECT 
    *
FROM trades.ohlcv_1m(
    start = '2024-09-28 18:00:00',
    end = '2024-09-28 19:00:00',
    base = 'bitcoin',
    quote = 'tether',
    exchange = 'gdax'
)
LIMIT 20

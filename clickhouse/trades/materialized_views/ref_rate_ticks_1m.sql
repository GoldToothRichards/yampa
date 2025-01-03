-- Table for 1m price ticks
CREATE TABLE IF NOT EXISTS trades.ref_rate_ticks_1m (
    base LowCardinality(String),
    quote LowCardinality(String),
    source LowCardinality(String),
    window_start DateTime,
    window_end DateTime,
    open_state AggregateFunction(argMin, Float64, DateTime),
    high_state AggregateFunction(max, Float64),
    low_state AggregateFunction(min, Float64),
    close_state AggregateFunction(argMax, Float64, DateTime),
    volume_base_state AggregateFunction(sum, Float64),
    volume_quote_state AggregateFunction(sum, Float64),
    trade_count_state AggregateFunction(count, UInt64)
)
ENGINE = AggregatingMergeTree()
PRIMARY KEY (base, quote, source)
ORDER BY (base, quote, source, window_start)
PARTITION BY toYYYYMM(window_start)
TTL
    window_start + INTERVAL 1 MONTH TO VOLUME 's3',
    window_start + INTERVAL 1 YEAR DELETE;

-- 1m ref rate ticks
CREATE MATERIALIZED VIEW trades.ref_rate_ticks_1m_mv TO trades.ref_rate_ticks_1m AS 
WITH clean_trades AS
(
    SELECT * FROM trades.raw_trades
    WHERE source != ''
      AND base != ''
      AND quote != ''
      AND price > 0
      AND volume_base > 0
      AND volume_quote > 0
      AND base != quote
)
SELECT
    base,
    quote,
    source,
    tumbleStart(toDateTime(timestamp), toIntervalMinute(1)) AS window_start,
    tumbleEnd(toDateTime(timestamp), toIntervalMinute(1)) AS window_end,
    argMinState(price, toDateTime(timestamp)) AS open_state,
    maxState(price) AS high_state,
    minState(price) AS low_state,
    argMaxState(price, toDateTime(timestamp)) AS close_state,
    sumState(volume_base) AS volume_base_state,
    sumState(volume_quote) AS volume_quote_state,
    countState() AS trade_count_state
FROM clean_trades
GROUP BY
    base,
    quote,
    source,
    window_start,
    window_end;

-- OHLCV View
CREATE OR REPLACE VIEW trades.ohlcv_ref_rate_1m AS
SELECT
    base,
    quote,
    source,
    window_start,
    argMinMerge(open_state) as open,
    maxMerge(high_state) as high,
    minMerge(low_state) as low,
    argMaxMerge(close_state) as close,
    sumMerge(volume_base_state) as volume_base,
    sumMerge(volume_quote_state) as volume_quote,
    countMerge(trade_count_state) as count,
    (close - open) / open as pct_change,
    volume_quote / volume_base as vwap
FROM trades.ref_rate_ticks_1m
WHERE (window_start >= {start: String}) AND (window_start <= {end: String})
GROUP BY
    base,
    quote,
    source,
    window_start
ORDER BY
    base,
    quote,
    source,
    window_start ASC
WITH FILL
    FROM toDateTime({start: String})
    TO toDateTime({end: String})
    STEP INTERVAL 1 MINUTE
    INTERPOLATE (vwap)

-- Test all pairs
SELECT 
    *
FROM trades.ohlcv_ref_rate_1m(
    start = '2024-09-28 18:00:00',
    end = '2024-09-28 19:00:00',
)
LIMIT 20

-- Test single pair
SELECT 
    *
FROM trades.ohlcv_ref_rate_1m(
    start = '2024-09-29 21:58:00',
    end = '2024-09-29 22:58:00'
)
WHERE
    base = 'bitcoin'
    AND quote = 'united-states-dollar'
    AND volume_base > 0
LIMIT 20


-- Test params
SET param_start = '2024-09-29 21:58:00';
SET param_end = '2024-09-29 22:58:00';
SET param_base = 'bitcoin';
SET param_quote = 'tether';
SET param_exchange = 'gdax';

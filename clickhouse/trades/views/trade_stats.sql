-- Stats with a parameterized time window for each distinct (base, quote, exchange)
CREATE OR REPLACE VIEW trades.exchange_stats_view AS
    SELECT
        base,
        quote,
        exchange,
        {start_time:DateTime} AS window_start,
        {end_time:DateTime} AS window_end,
        count() AS trade_count,
        -- Volume Base stats
        min(volume_base) AS min_volume_base,
        max(volume_base) AS max_volume_base,
        avg(volume_base) AS avg_volume_base,
        median(volume_base) AS median_volume_base,
        sum(volume_base) AS sum_volume_base,
        stddevPop(volume_base) AS stddev_volume_base,
        -- Volume Quote stats
        min(volume_quote) AS min_volume_quote,
        max(volume_quote) AS max_volume_quote,
        avg(volume_quote) AS avg_volume_quote,
        median(volume_quote) AS median_volume_quote,
        sum(volume_quote) AS sum_volume_quote,
        stddevPop(volume_quote) AS stddev_volume_quote,
        -- Price stats
        min(price) AS min_price,
        max(price) AS max_price,
        avg(price) AS avg_price,
        median(price) AS median_price,
        stddevPop(price) AS stddev_price,
        sum(volume_quote) / sum(volume_base) AS vwap
    FROM trades.clean_trades_view
    WHERE timestamp >= {start_time:DateTime} AND timestamp < {end_time:DateTime}
    GROUP BY base, quote, exchange;

-- Stats with a parameterized time window for each distinct (base, quote)
-- Pair stats are aggregated over all exchanges
CREATE OR REPLACE VIEW trades.pair_stats_view AS
    SELECT
        base,
        quote,
        {start_time:DateTime} AS window_start,
        {end_time:DateTime} AS window_end,
        count() AS trade_count,
        -- Volume Base stats
        min(volume_base) AS min_volume_base,
        max(volume_base) AS max_volume_base,
        avg(volume_base) AS avg_volume_base,
        median(volume_base) AS median_volume_base,
        sum(volume_base) AS sum_volume_base,
        stddevPop(volume_base) AS stddev_volume_base,
        -- Volume Quote stats
        min(volume_quote) AS min_volume_quote,
        max(volume_quote) AS max_volume_quote,
        avg(volume_quote) AS avg_volume_quote,
        median(volume_quote) AS median_volume_quote,
        sum(volume_quote) AS sum_volume_quote,
        stddevPop(volume_quote) AS stddev_volume_quote,
        -- Price stats
        min(price) AS min_price,
        max(price) AS max_price,
        avg(price) AS avg_price,
        median(price) AS median_price,
        stddevPop(price) AS stddev_price,
        sum(volume_quote) / sum(volume_base) AS vwap
    FROM trades.clean_trades_view
    WHERE timestamp >= {start_time:DateTime} AND timestamp < {end_time:DateTime}
    GROUP BY base, quote;

-- Query for stats aggregated over past 24 hours
SELECT *
FROM trades.exchange_stats_view(
    start_time = subtractHours(now(), 24),
    end_time = now()
)
LIMIT 50;

SELECT *
FROM trades.pair_stats_view(
    start_time = subtractHours(now(), 24),
    end_time = now()
)
LIMIT 50;

-- View for outliers
CREATE OR REPLACE VIEW trades.outliers_view AS
WITH stats AS (
    SELECT
        base,
        quote,
        exchange,
        timestamp,
        price,
        volume_base,
        volume_quote,
        avg(price) OVER w AS avg_price,
        stddevPop(price) OVER w AS stddev_price,
        avg(volume_base) OVER w AS avg_volume,
        stddevPop(volume_base) OVER w AS stddev_volume,
        {threshold:Float64} AS threshold,
        {start_time:DateTime} AS window_start,
        {end_time:DateTime} AS window_end
    FROM trades.clean_trades_view
    WHERE timestamp >= {start_time:DateTime} AND timestamp < {end_time:DateTime}
    WINDOW w AS (PARTITION BY base, quote)
)
SELECT *
FROM stats
WHERE
    (price > avg_price + threshold * stddev_price) OR
    (price < avg_price - threshold * stddev_price) OR
    (volume_base > avg_volume + threshold * stddev_volume) OR
    (volume_base < avg_volume - threshold * stddev_volume);

-- View for inliers
CREATE OR REPLACE VIEW trades.inliers_view AS
WITH stats AS (
    SELECT
        base,
        quote,
        exchange,
        timestamp,
        price,
        volume_base,
        volume_quote,
        avg(price) OVER w AS avg_price,
        stddevPop(price) OVER w AS stddev_price,
        avg(volume_base) OVER w AS avg_volume,
        stddevPop(volume_base) OVER w AS stddev_volume,
        {threshold:Float64} AS threshold,
        {start_time:DateTime} AS window_start,
        {end_time:DateTime} AS window_end
    FROM trades.clean_trades_view
    WHERE timestamp >= {start_time:DateTime} AND timestamp < {end_time:DateTime}
    WINDOW w AS (PARTITION BY base, quote)
)
SELECT *
FROM stats
WHERE
    (price <= avg_price + threshold * stddev_price) AND
    (price >= avg_price - threshold * stddev_price) AND
    (volume_base <= avg_volume + threshold * stddev_volume) AND
    (volume_base >= avg_volume - threshold * stddev_volume);

-- Example query to use the views
SELECT *
FROM trades.inliers_view(
    threshold = 2,
    start_time = subtractHours(now(), 24),
    end_time = now()
)
LIMIT 50;

SELECT *
FROM trades.outliers_view(
    threshold = 2,
    start_time = subtractHours(now(), 24),
    end_time = now()
)
LIMIT 50;


-- Compare counts
SELECT 
    'raw_trades' AS category,
    COUNT(*) AS count
FROM trades.raw_trades
WHERE timestamp >= subtractHours(now(), 24) AND timestamp < now()

UNION ALL

SELECT 
    'clean_trades' AS category,
    COUNT(*) AS count
FROM trades.clean_trades_view
WHERE timestamp >= subtractHours(now(), 24) AND timestamp < now()

UNION ALL

SELECT 
    'rejected_trades' AS category,
    COUNT(*) AS count
FROM trades.rejected_trades_view
WHERE timestamp >= subtractHours(now(), 24) AND timestamp < now()

UNION ALL

SELECT 
    'inliers' AS category,
    COUNT(*) AS count
FROM trades.inliers_view(
    threshold = 2,
    start_time = subtractHours(now(), 24),
    end_time = now()
)

UNION ALL

SELECT 
    'outliers' AS category,
    COUNT(*) AS count
FROM trades.outliers_view(
    threshold = 2,
    start_time = subtractHours(now(), 24),
    end_time = now()
);

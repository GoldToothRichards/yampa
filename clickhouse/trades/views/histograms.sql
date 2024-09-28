-- View for volume histogram
CREATE OR REPLACE VIEW trades.volume_histogram_view AS
WITH params AS (
    SELECT
        {base:String} AS base,
        {quote:String} AS quote,
        {num_bins:UInt32} AS num_bins,
        {start_time:DateTime} AS start_time,
        {end_time:DateTime} AS end_time,
        {threshold:Float64} AS threshold
),
inliers AS (
    SELECT *
    FROM trades.inliers_view(
        threshold = (SELECT threshold FROM params),
        start_time = (SELECT start_time FROM params),
        end_time = (SELECT end_time FROM params)
    )
    WHERE base = (SELECT base FROM params) AND quote = (SELECT quote FROM params)
),
stats AS (
    SELECT
        min(volume_base) AS min_volume,
        max(volume_base) AS max_volume
    FROM inliers
),
histogram AS (
    SELECT
        bin,
        bin_start,
        bin_end,
        count() AS count
    FROM (
        SELECT
            volume_base,
            arrayJoin(range(1, num_bins + 1)) AS bin,
            min_volume + (bin - 1) * (max_volume - min_volume) / num_bins AS bin_start,
            min_volume + bin * (max_volume - min_volume) / num_bins AS bin_end
        FROM inliers, stats, params
    )
    WHERE volume_base >= bin_start AND volume_base < bin_end
    GROUP BY bin, bin_start, bin_end
    ORDER BY bin
)
SELECT
    bin,
    bin_start,
    bin_end,
    count
FROM histogram;

-- View for price histogram
CREATE OR REPLACE VIEW trades.price_histogram_view AS
WITH params AS (
    SELECT
        {base:String} AS base,
        {quote:String} AS quote,
        {num_bins:UInt32} AS num_bins,
        {start_time:DateTime} AS start_time,
        {end_time:DateTime} AS end_time,
        {threshold:Float64} AS threshold
),
inliers AS (
    SELECT *
    FROM trades.inliers_view(
        threshold = (SELECT threshold FROM params),
        start_time = (SELECT start_time FROM params),
        end_time = (SELECT end_time FROM params)
    )
    WHERE base = (SELECT base FROM params) AND quote = (SELECT quote FROM params)
),
stats AS (
    SELECT
        min(price) AS min_price,
        max(price) AS max_price
    FROM inliers
),
histogram AS (
    SELECT
        bin,
        bin_start,
        bin_end,
        count() AS count
    FROM (
        SELECT
            price,
            arrayJoin(range(1, num_bins + 1)) AS bin,
            min_price + (bin - 1) * (max_price - min_price) / num_bins AS bin_start,
            min_price + bin * (max_price - min_price) / num_bins AS bin_end
        FROM inliers, stats, params
    )
    WHERE price >= bin_start AND price < bin_end
    GROUP BY bin, bin_start, bin_end
    ORDER BY bin
)
SELECT
    bin,
    bin_start,
    bin_end,
    count
FROM histogram;

-- Query the volume histogram
SELECT *
FROM trades.volume_histogram_view(
    base = 'bitcoin',
    quote = 'tether',
    num_bins = 600,
    start_time = subtractHours(now(), 12),
    end_time = now(),
    threshold = 2
);

-- Query the price histogram
SELECT *
FROM trades.price_histogram_view(
    base = 'bitcoin',
    quote = 'tether',
    num_bins = 600,
    start_time = subtractHours(now(), 24),
    end_time = now(),
    threshold = 2
);
INSERT INTO rejected_trades_sink
SELECT * FROM trades_source
WHERE
    source IS NULL OR source = '' OR
    base IS NULL OR base = '' OR
    quote IS NULL OR quote = '' OR
    exchange IS NULL OR exchange = '' OR
    timestamp IS NULL OR
    price IS NULL OR price <= 0 OR
    volume_base IS NULL OR volume_base <= 0 OR
    volume_quote IS NULL OR volume_quote <= 0 OR
    base = quote;
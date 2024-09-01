INSERT INTO clean_trades_sink
SELECT * FROM trades_source
WHERE
    source IS NOT NULL AND source != '' AND
    base IS NOT NULL AND base != '' AND
    quote IS NOT NULL AND quote != '' AND
    exchange IS NOT NULL AND exchange != '' AND
    timestamp IS NOT NULL AND
    price IS NOT NULL AND price > 0 AND
    volume_base IS NOT NULL AND volume_base > 0 AND
    volume_quote IS NOT NULL AND volume_quote > 0;

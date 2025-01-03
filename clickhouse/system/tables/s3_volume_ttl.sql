ALTER TABLE system.processors_profile_log
MODIFY SETTING storage_policy = 'shared';

ALTER TABLE system.processors_profile_log
MODIFY TTL event_date + INTERVAL 1 DAY TO VOLUME 's3';
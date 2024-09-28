-- Settings changed
SELECT *
FROM system.settings
WHERE changed;

-- Inspect Storage Policies
SELECT
    policy_name,
    volume_name,
    volume_priority,
    disks
FROM system.storage_policies


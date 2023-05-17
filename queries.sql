--name: set-cache
-- Save a phone number with version to the cache
-- $1: phone_number
-- $2: version
INSERT INTO segmentation(phone_number, canary_version) VALUES($1, $2) RETURNING id

--name: get-cache
-- Load phone number canary version
-- $1: phone_number
SELECT canary_version FROM segmentation WHERE phone_number=$1

--name: update-cache
-- Update phone number canary version
-- $1: phone_number
-- $2: canary_version
UPDATE segmentation SET canary_version = $2 WHERE phone_number=$1
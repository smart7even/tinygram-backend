-- extend device id and device token size to 255 characters
ALTER TABLE devices
    ALTER COLUMN device_id TYPE VARCHAR(255),
    ALTER COLUMN device_token TYPE VARCHAR(255);
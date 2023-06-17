CREATE TABLE public.device (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  device_id STRING(20) NOT NULL,
  connection_status STRING(20) NULL,
  device_type INT NOT_NULL,
  firmware_version STRING(20) NULL,
  monitoring_status STRING(20) NULL,
  nickname STRING(40) NULL,
  temperature INT NULL,
  "owner" UUID NULL,
  config UUID NULL,
  temperature_timeline UUID NULL,
  last_seen TIMESTAMPTZ NULL,
)
CREATE TABLE public.device (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  connection_status INT NULL,
  device_type INT NOT_NULL,
  firmware_version STRING(30) NULL,
  monitoring_status INT NULL,
  nickname STRING(40) NULL,
  temperature INT NULL,
  "owner" UUID NULL,
  last_seen TIMESTAMPTZ NULL,
)
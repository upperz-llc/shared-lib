CREATE TABLE public.device (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  connection_status INT NOT_NULL,
  device_type INT NOT_NULL,
  monitoring_status INT NULL,
  firmware_version STRING(30) NULL,
  "owner" UUID NULL,
  created_at TIMESTAMPTZ DEFAULT now() NOT_NULL,
  updated_at TIMESTAMPTZ NULL
)
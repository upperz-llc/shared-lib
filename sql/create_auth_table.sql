CREATE TABLE public.auth (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  device_id UUID NOT NULL,
  username STRING(60) NULL,
  password STRING(120) NULL,
  enabled BOOL DEFAULT true NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
  INDEX (device_id)
)
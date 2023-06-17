CREATE TABLE public.alarm (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "type" INT NOT NULL,
  device UUID NOT NULL,
  device_id STRING(20) NOT NULL,
  created_at TIMESTAMPTZ NULL,
  closed_at TIMESTAMPTZ NULL,
  acked_at TIMESTAMPTZ NULL,
  acked_by UUID NULL,
  acked_check_count INT NULL,
  acked BOOL NULL,
  active BOOL NULL
)
CREATE TABLE public.device (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  device_type INT NOT NULL,
  measurement_type INT NOT NULL,
  connection_status INT NOT NULL,
  monitoring_status INT NULL,
  firmware_version STRING(30) NULL,
  "owner" UUID NULL REFERENCES public.user(id),
  last_seen_at TIMSTAMPTZ NULL,
  created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
  updated_at TIMESTAMPTZ NULL
)
CREATE TABLE public.device_config (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  device_id UUID NOT NULL REFERENCES public.device(id),
  alert INT NULL,
  warning INT NULL,
  'target' INT NULL,
  measurement_interval INT NULL,
  'version' INT NULL,
  created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
  updated_at TIMESTAMPTZ NULL
)
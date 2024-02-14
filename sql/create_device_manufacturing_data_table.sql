CREATE TABLE public.device_manufacturing_data (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  device_id UUID NOT NULL REFERENCES public.device(id),
  device_type INT NOT NULL,
  measurement_type INT NULL,
  username STRING(60) NULL,
  password STRING(120) NULL,
  manufactured_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ DEFAULT now() NOT NULL
)
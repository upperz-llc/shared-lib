CREATE TABLE public.device_manufacturing_info (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  device_id STRING(20) NOT NULL,
  device_type INT NULL,
  manufactured_date TIMESTAMPTZ NULL,
  manufactured_type INT NULL,
  username STRING(60) NULL,
)
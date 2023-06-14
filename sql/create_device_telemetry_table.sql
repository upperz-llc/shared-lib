CREATE TABLE public.temperature (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  temperature DECIMAL(20) NULL,
  device_id STRING(20) NULL,
  "timestamp" TIMESTAMPTZ NULL,
  CONSTRAINT temperature_pkey PRIMARY KEY (id ASC),
  INDEX temperature_device_id_timestamp_idx (device_id ASC, "timestamp" ASC) STORING (temperature)
)
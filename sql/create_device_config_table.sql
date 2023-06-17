CREATE TABLE public.device_config (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  alert_temperature INT NULL,
  warning_temperature INT NULL,
  target_temperature INT NULL,
  telemetry_interval int NULL,
)
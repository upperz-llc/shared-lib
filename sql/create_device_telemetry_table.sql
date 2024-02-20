CREATE TABLE public.telemetry (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  device_id UUID NOT NULL REFERENCES public.device(id),
  "value" DECIMAL(20) NOT NULL,
  "type" INT NOT NULL,
  "timestamp" TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
)
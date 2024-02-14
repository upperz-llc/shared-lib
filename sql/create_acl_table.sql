CREATE TABLE public.acl (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  device_id UUID NOT NULL,
  auth_id UUID NOT NULL,
  topic STRING(240) NULL,
  access STRING(6) NULL,
  allowed BOOL DEFAULT true NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
  INDEX (device_id, topic)
)
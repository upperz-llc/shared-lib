CREATE TABLE public.user (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  uid STRING(40) NOT NULL,
  email STRING(100) NULL,
  notification_push BOOL NULL,
  notification_sms BOOL NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NULL
  )




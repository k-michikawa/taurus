CREATE TABLE public.m_users (
  id uuid NOT NULL DEFAULT uuid_generate_v4(),
  name text NOT NULL,
  email text NOT NULL,
  password text NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp,
  is_deleted boolean NOT NULL DEFAULT 'f',
  CONSTRAINT m_users_pk PRIMARY KEY (id)
);

GRANT INSERT, SELECT, UPDATE, DELETE ON public.m_users TO taurus;

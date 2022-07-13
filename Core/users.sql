CREATE TABLE public.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    balance FLOAT NOT NULL
);
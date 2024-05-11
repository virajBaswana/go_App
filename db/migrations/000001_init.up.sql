CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    first_name character varying(100) NOT NULL,
    last_name character varying(100) NOT NULL,
    email character varying(100) UNIQUE NOT NULL,
    password character varying(100) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS connections (
    id serial PRIMARY KEY,
    target_id integer REFERENCES users (id),
    initiator_id integer REFERENCES users (id),
    is_reciprocated boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    UNIQUE(target_id , initiator_id)
);
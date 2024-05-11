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
    target_id integer REFERENCES users (id),
    initiator_id integer REFERENCES users (id),
    is_reciprocated boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
    PRIMARY KEY(target_id , initiator_id)
);

INSERT INTO
    users (first_name, last_name, email, password)
VALUES
    (
        'tommy',
        'robinson',
        'tomrobins@gmail.com',
        '$2a$10$.6fPhbScH.ygrm6v6BSkA.bOo90nlMEMDS/IdlcBEzDRi759lwYFu'
    ),
    (
       'dick',
       'green',
       'greendicky@gmail.com',
       '$2a$10$.6fPhbScH.ygrm6v6BSkA.bOo90nlMEMDS/IdlcBEzDRi759lwYFu'
    ),
    (
        'harry',
        'potter',
        'harrypotter@gmail.com',
        '$2a$10$.6fPhbScH.ygrm6v6BSkA.bOo90nlMEMDS/IdlcBEzDRi759lwYFu'
    ) ON CONFLICT (email) DO NOTHING;
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

INSERT INTO
    connections(target_id, initiator_id, is_reciprocated)
VALUES
    (1, 2, false),
    (1, 3, false),
    (2, 3, true) ON CONFLICT (target_id, initiator_id) DO NOTHING
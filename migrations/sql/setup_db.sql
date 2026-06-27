INSERT INTO appuser (username) VALUES
    ('test_user_1'),
    ('test_user_2'),
    ('test_user_3'),
    ('test_user_4'),
    ('test_user_5');

INSERT INTO posts (user_id, title, text) VALUES
    (1, 'post 1',  'simple text 1'),
    (2, 'post 2',  'simple text 2'),
    (3, 'post 3',  'simple text 3'),
    (4, 'post 4',  'simple text 4'),
    (5, 'post 5',  'simple text 5'),
    (1, 'post 6',  'simple text 6'),
    (2, 'post 7',  'simple text 7'),
    (3, 'post 8',  'simple text 8'),
    (4, 'post 9',  'simple text 9'),
    (5, 'post 10', 'simple text 10');
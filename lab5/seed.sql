CREATE EXTENSION IF NOT EXISTS pgcrypto;

BEGIN;

INSERT INTO teachers (
    last_name,
    first_name,
    middle_name,
    login,
    password_hash
)
VALUES
(
    'Иванов',
    'Иван',
    'Иванович',
    'teacher1',
    crypt('12345', gen_salt('bf'))
),
(
    'Петров',
    'Пётр',
    'Петрович',
    'teacher2',
    crypt('12345', gen_salt('bf'))
)
ON CONFLICT (login) DO NOTHING;

INSERT INTO student_groups (
    group_name,
    group_prefix,
    group_number
)
VALUES
('ИВ-222', 'ИВ', '222'),
('ИВ-223', 'ИВ', '223'),
('ИС-221', 'ИС', '221'),
('ИС-222', 'ИС', '222'),
('ИП-201', 'ИП', '201'),
('ИП-202', 'ИП', '202')
ON CONFLICT (group_name) DO NOTHING;

WITH g AS (
    SELECT group_uuid
    FROM student_groups
    WHERE group_name = 'ИВ-222'
)
INSERT INTO students (
    last_name,
    first_name,
    middle_name,
    login,
    password_hash,
    group_uuid
)
SELECT
    s.last_name,
    s.first_name,
    s.middle_name,
    s.login,
    crypt('12345', gen_salt('bf')),
    g.group_uuid
FROM (
    VALUES
        ('Сидоров', 'Сергей', 'Сергеевич', 'student1'),
        ('Кузнецов', 'Кирилл', 'Кириллович', 'student2')
) AS s(last_name, first_name, middle_name, login)
CROSS JOIN g
ON CONFLICT (login) DO NOTHING;

WITH g AS (
    SELECT group_uuid
    FROM student_groups
    WHERE group_name = 'ИВ-223'
)
INSERT INTO students (
    last_name,
    first_name,
    middle_name,
    login,
    password_hash,
    group_uuid
)
SELECT
    s.last_name,
    s.first_name,
    s.middle_name,
    s.login,
    crypt('12345', gen_salt('bf')),
    g.group_uuid
FROM (
    VALUES
        ('Смирнов', 'Алексей', 'Игоревич', 'student3')
) AS s(last_name, first_name, middle_name, login)
CROSS JOIN g
ON CONFLICT (login) DO NOTHING;

WITH g AS (
    SELECT group_uuid
    FROM student_groups
    WHERE group_name = 'ИС-221'
)
INSERT INTO students (
    last_name,
    first_name,
    middle_name,
    login,
    password_hash,
    group_uuid
)
SELECT
    s.last_name,
    s.first_name,
    s.middle_name,
    s.login,
    crypt('12345', gen_salt('bf')),
    g.group_uuid
FROM (
    VALUES
        ('Орлов', 'Дмитрий', 'Андреевич', 'student4'),
        ('Морозов', 'Никита', 'Владимирович', 'student5')
) AS s(last_name, first_name, middle_name, login)
CROSS JOIN g
ON CONFLICT (login) DO NOTHING;

WITH g AS (
    SELECT group_uuid
    FROM student_groups
    WHERE group_name = 'ИС-222'
)
INSERT INTO students (
    last_name,
    first_name,
    middle_name,
    login,
    password_hash,
    group_uuid
)
SELECT
    s.last_name,
    s.first_name,
    s.middle_name,
    s.login,
    crypt('12345', gen_salt('bf')),
    g.group_uuid
FROM (
    VALUES
        ('Попов', 'Егор', 'Максимович', 'student6')
) AS s(last_name, first_name, middle_name, login)
CROSS JOIN g
ON CONFLICT (login) DO NOTHING;

WITH g AS (
    SELECT group_uuid
    FROM student_groups
    WHERE group_name = 'ИП-201'
)
INSERT INTO students (
    last_name,
    first_name,
    middle_name,
    login,
    password_hash,
    group_uuid
)
SELECT
    s.last_name,
    s.first_name,
    s.middle_name,
    s.login,
    crypt('12345', gen_salt('bf')),
    g.group_uuid
FROM (
    VALUES
        ('Васильев', 'Роман', 'Олегович', 'student7')
) AS s(last_name, first_name, middle_name, login)
CROSS JOIN g
ON CONFLICT (login) DO NOTHING;

WITH g AS (
    SELECT group_uuid
    FROM student_groups
    WHERE group_name = 'ИП-202'
)
INSERT INTO students (
    last_name,
    first_name,
    middle_name,
    login,
    password_hash,
    group_uuid
)
SELECT
    s.last_name,
    s.first_name,
    s.middle_name,
    s.login,
    crypt('12345', gen_salt('bf')),
    g.group_uuid
FROM (
    VALUES
        ('Фёдоров', 'Артём', 'Ильич', 'student8')
) AS s(last_name, first_name, middle_name, login)
CROSS JOIN g
ON CONFLICT (login) DO NOTHING;

COMMIT;

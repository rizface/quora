-- password: testdata
TRUNCATE accounts CASCADE;

INSERT INTO accounts(id, email, username, password) VALUES (
    'f028ac5a-e4c9-442f-bf9a-86c024a79baa', 'testlogin@gmail.com', 'testlogin', '$2a$10$NXt0Uyr7XIOEM4TQZMVd7uWVkuTcG8pqTsBFZHAGg86.jjp32VtjW'
);

INSERT INTO spaces(id, owner_id, name) VALUES
('a53152d7-2d24-42e1-a55f-649e87349ffa', 'f028ac5a-e4c9-442f-bf9a-86c024a79baa', 'Ruang Programmer');
-- password: testdata
TRUNCATE accounts CASCADE;

INSERT INTO accounts(id, email, username, password) VALUES (
    'f028ac5a-e4c9-442f-bf9a-86c024a79baa', 'testlogin@gmail.com', 'testlogin', '$2a$10$NXt0Uyr7XIOEM4TQZMVd7uWVkuTcG8pqTsBFZHAGg86.jjp32VtjW'
);

INSERT INTO spaces(id, owner_id, name) VALUES
('a53152d7-2d24-42e1-a55f-649e87349ffa', 'f028ac5a-e4c9-442f-bf9a-86c024a79baa', 'Ruang Programmer');

INSERT INTO public.questions (id,author_id,space_id,question,upvote,downvote,created_at,updated_at) VALUES
('4b9ef364-0d6a-4f60-a169-39b1d076c65d','f028ac5a-e4c9-442f-bf9a-86c024a79baa',NULL,'yooo this work >?',0,0,'2023-09-02 02:42:59.334677','2023-09-02 02:42:59.334677'),
('4b9ef364-0d6a-4f60-a169-39b1d076c65e','f028ac5a-e4c9-442f-bf9a-86c024a79baa',NULL,'yooo this work >?',1,0,'2023-09-02 02:42:59.334677','2023-09-02 02:42:59.334677'),
('4b9ef364-0d6a-4f60-a169-39b1d076c65f','f028ac5a-e4c9-442f-bf9a-86c024a79baa',NULL,'yooo this work >?',0,0,'2023-09-02 02:42:59.334677','2023-09-02 02:42:59.334677'),
('4b9ef364-0d6a-4f60-a169-39b1d076c65b','f028ac5a-e4c9-442f-bf9a-86c024a79baa',NULL,'yooo this work >?',0,1,'2023-09-02 02:42:59.334677','2023-09-02 02:42:59.334677');

INSERT INTO votes(voter_id, question_id, "type") VALUES 
('f028ac5a-e4c9-442f-bf9a-86c024a79baa', '4b9ef364-0d6a-4f60-a169-39b1d076c65e', 'upvote'),
('f028ac5a-e4c9-442f-bf9a-86c024a79baa', '4b9ef364-0d6a-4f60-a169-39b1d076c65b', 'downvote');
-- password: testdata
TRUNCATE accounts CASCADE;

INSERT INTO accounts(id, email, username, password) VALUES (
    'f028ac5a-e4c9-442f-bf9a-86c024a79baa', 'testlogin@gmail.com', 'testlogin', '$2a$10$NXt0Uyr7XIOEM4TQZMVd7uWVkuTcG8pqTsBFZHAGg86.jjp32VtjW'
);

INSERT INTO spaces(id, owner_id, name) VALUES
('a53152d7-2d24-42e1-a55f-649e87349ffa', 'f028ac5a-e4c9-442f-bf9a-86c024a79baa', 'Ruang Programmer');

INSERT INTO questions (id,author_id,space_id,question,created_at,updated_at) VALUES
('4b9ef364-0d6a-4f60-a169-39b1d076c65d','f028ac5a-e4c9-442f-bf9a-86c024a79baa',NULL,'yooo this work >?','2023-09-02 02:42:59.334677','2023-09-02 02:42:59.334677'),
('4b9ef364-0d6a-4f60-a169-39b1d076c65e','f028ac5a-e4c9-442f-bf9a-86c024a79baa',NULL,'yooo this work >?','2023-09-02 02:42:59.334677','2023-09-02 02:42:59.334677'),
('4b9ef364-0d6a-4f60-a169-39b1d076c65f','f028ac5a-e4c9-442f-bf9a-86c024a79baa',NULL,'yooo this work >?','2023-09-02 02:42:59.334677','2023-09-02 02:42:59.334677'),
('4b9ef364-0d6a-4f60-a169-39b1d076c65b','f028ac5a-e4c9-442f-bf9a-86c024a79baa',NULL,'yooo this work >?','2023-09-02 02:42:59.334677','2023-09-02 02:42:59.334677');

INSERT INTO answers(id, question_id, answerer_id, upvote, downvote, answer) VALUES
('4b9ef364-0d6a-4f60-a169-39b1d076c65d', '4b9ef364-0d6a-4f60-a169-39b1d076c65e', 'f028ac5a-e4c9-442f-bf9a-86c024a79baa',0,0,'answer 1'),
('4b9ef364-0d6a-4f60-a169-39b1d076c65f', '4b9ef364-0d6a-4f60-a169-39b1d076c65e', 'f028ac5a-e4c9-442f-bf9a-86c024a79baa',0,0,'answer 2'),
('4b9ef364-0d6a-4f60-a169-39b1d076c65e', '4b9ef364-0d6a-4f60-a169-39b1d076c65e', 'f028ac5a-e4c9-442f-bf9a-86c024a79baa',1,0,'answer 3'),
('4b9ef364-0d6a-4f60-a169-39b1d076c65b', '4b9ef364-0d6a-4f60-a169-39b1d076c65e', 'f028ac5a-e4c9-442f-bf9a-86c024a79baa',0,1,'answer 4');


INSERT INTO votes(voter_id, answer_id, "type") VALUES 
('f028ac5a-e4c9-442f-bf9a-86c024a79baa', '4b9ef364-0d6a-4f60-a169-39b1d076c65e', 'upvote'),
('f028ac5a-e4c9-442f-bf9a-86c024a79baa', '4b9ef364-0d6a-4f60-a169-39b1d076c65b', 'downvote');
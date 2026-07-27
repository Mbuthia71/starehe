INSERT INTO groups (id, name, type, join_policy, description) 
VALUES (gen_random_uuid(), 'Test Chapter', 'chapter', 'open', 'Test group for presentation demo') 
RETURNING id;

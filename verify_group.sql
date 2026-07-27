SELECT g.id, g.name, g.type, gm.role 
FROM groups g 
LEFT JOIN group_members gm ON g.id = gm.group_id 
WHERE g.id = '9ee34e8b-1b12-4e86-b6ab-4c3c249d02f0';

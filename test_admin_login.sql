SELECT id, email, role FROM users WHERE role = 'super_admin' OR role = 'admin' LIMIT 5;

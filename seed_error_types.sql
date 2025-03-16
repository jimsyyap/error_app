INSERT INTO Error_Types (name) VALUES 
('Forehand'),
('Backhand'),
('Serve'),
('Volley')
ON CONFLICT (name) DO NOTHING;

-- Insert sample movies
INSERT INTO movies (title, description, duration_minutes, release_date, created_at)
VALUES 
('Spider-Man: No Way Home', 'Peter Parker mengatasi konsekuensi dari identitas Spider-Man yang terungkap', 148, '2021-12-15', NOW()),
('Avatar: The Way of Water', 'Jake Sully dan keluarganya berjuang untuk tetap bersama', 192, '2022-12-14', NOW()),
('Top Gun: Maverick', 'Pilot test yang berani menghadapi masa lalu dan masa depan', 130, '2022-05-24', NOW());

-- Insert sample cinemas
INSERT INTO cinemas (name, city, address, created_at)
VALUES 
('MKP XXI Jakarta Pusat', 'Jakarta', 'Jl. Thamrin No. 1', NOW()),
('MKP XXI Bandung', 'Bandung', 'Jl. Asia Afrika No. 10', NOW());

-- Insert sample studios
INSERT INTO studios (cinema_id, name, total_seats, created_at)
VALUES 
(1, 'Studio 1', 100, NOW()),
(1, 'Studio 2', 80, NOW()),
(1, 'IMAX', 150, NOW()),
(2, 'Studio 1', 100, NOW());

-- Insert sample schedules
INSERT INTO schedules (movie_id, studio_id, start_time, end_time, price, status, created_at)
VALUES 
(1, 1, '2024-12-05 14:00:00', '2024-12-05 16:30:00', 50000, 'SHOWING', NOW()),
(1, 1, '2024-12-05 19:00:00', '2024-12-05 21:30:00', 50000, 'SHOWING', NOW()),
(2, 3, '2024-12-05 15:00:00', '2024-12-05 18:15:00', 75000, 'SHOWING', NOW()),
(3, 2, '2024-12-05 16:00:00', '2024-12-05 18:15:00', 45000, 'SHOWING', NOW());
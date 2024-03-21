CREATE TABLE access_control (
    id SERIAL PRIMARY KEY,
    role VARCHAR(50) NOT NULL,
    endpoint VARCHAR(100) NOT NULL,
    method VARCHAR(10) NOT NULL
);

INSERT INTO access_control (role, endpoint, method) VALUES
('unauthorized', '/v1/verify/', 'GET'),
('unauthorized', '/v1/register/', 'POST'),
('unauthorized', '/v1/login/', 'POST'),
('unauthorized', '/swagger/*', 'GET'),
('user', '/v1/getall/', 'GET'),
('user', '/v1/delete/', 'DELETE'),
('user', '/v1/update/', 'PUT'),
('user', '/v1/getuser/', 'GET');
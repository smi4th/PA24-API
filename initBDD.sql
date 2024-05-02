INSERT INTO ACCOUNT_TYPE (uuid, type, private, admin) VALUES
('1', 'Locataire', 'false', 'false'),
('2', 'Loueur', 'false', 'false'),
('3', 'Administateur', 'true', 'true'),
('4', 'Provider', 'true', 'false'),
('5', 'Handyman', 'false', 'false');

INSERT INTO ACCOUNT (uuid, token, username, password, first_name, last_name, email, account_type) VALUES
('1', 'token1', 'user1', '$2a$10$8gd1puAaAsl4LBpvAj6uhO4whmoaSKiD69AiZGEamUtDOr6ALZVmG', 'John', 'Doe', 'john.doe@example.com', '1'),
('2', 'token2', 'user2', '$2a$10$rWQKIlPsWD0NIYVjV8xuROlL3ZfZu8nMVq.mIoHuz0fvlYS7aZBda', 'Jane', 'Smith', 'jane.smith@example.com', '2'),
('3', 'token3', 'user3', '$2a$10$Ooru.yReQwZ4v3XoTvQFgu.IgwyuXnaeEqBiVjrdJbWgLlri0juL.', 'Alice', 'Johnson', 'alice.johnson@example.com', '3'),
('4', 'token4', 'user4', '$2a$10$ZPgHF8EJ93a5JLDKCaEmJOhqde6CXepP/NZXhKP1EEbx3kK1BMZgm', 'Bob', 'Brown', 'bob.brown@example.com', '4'),
('5', 'token5', 'user5', '$2a$10$MDVKuspmpvdW/nYaYxpKKe94kYX1bpIj9u40E8GDgOEjBiv/gq5ne', 'Emma', 'Wilson', 'emma.wilson@example.com', '5');

INSERT INTO PROVIDER (uuid, name, email) VALUES
('1', 'Provider1', 'provider1@provider1.com'),
('2', 'Provider2', 'provider2@provider2.com'),
('3', 'Provider3', 'provider3@provider3.com');

INSERT INTO PROVIDER_ACCOUNT (administration_level, provider, account) VALUES
(1, '1', '4'),
(2, '2', '5');

INSERT INTO SUBSCRIPTION (uuid, name) VALUES
('1', 'Subscription1'),
('2', 'Subscription2'),
('3', 'Subscription3');

INSERT INTO ACCOUNT_SUBSCRIPTION (start_date, account, subscription) VALUES
(NOW(), '1', '1'),
(NOW(), '2', '2');

INSERT INTO SERVICES_TYPES (uuid, type) VALUES
('1', 'Service1'),
('2', 'Service2'),
('3', 'Service3');

INSERT INTO SERVICES (uuid, price, service_type) VALUES
('1', 10.00, '1'),
('2', 20.00, '2'),
('3', 30.00, '3');

INSERT INTO ACCOUNT_SERVICES (account, services) VALUES
('5', '1'),
('4', '2'),
('2', '3');

INSERT INTO CONSUME (report, notice, price, note, services, account) VALUES
('Report1', 'Notice1', 10.00, 4, '1', '1'),
('Report2', 'Notice2', 20.00, 5, '2', '1'),
('Report3', 'Notice3', 30.00, 3, '3', '1');

INSERT INTO DISPONIBILITY (uuid, start_date, end_date) VALUES
('1', NOW(), DATE_ADD(NOW(), INTERVAL 7 DAY)),
('2', NOW(), DATE_ADD(NOW(), INTERVAL 14 DAY)),
('3', NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY));

INSERT INTO DISPONIBILITY_ACCOUNT (disponibility, account) VALUES
('1', '1'),
('2', '2'),
('3', '3');

INSERT INTO HOUSE_TYPE (uuid, type) VALUES
('1', 'Maison individuelle'),
('2', 'Appartement'),
('3', 'Maison de ville');

INSERT INTO HOUSING (uuid, surface, price, validated, street_nb, city, zip_code, street, description, house_type) VALUES
('1', 100.00, 200000.00, true, '123', 'Paris', '75001', 'Rue de Rivoli', 'Belle maison individuelle', '1'),
('2', 75.50, 150000.00, false, '456', 'Lyon', '69001', 'Rue de la République', 'Appartement lumineux', '2'),
('3', 120.75, 300000.00, true, '789', 'Marseille', '13001', 'Avenue du Prado', 'Maison de ville avec jardin', '3');

INSERT INTO EQUIPMENT_TYPE (uuid, name) VALUES
('1', 'Électroménager'),
('2', 'Mobilier'),
('3', 'Électronique');

INSERT INTO EQUIPMENT (uuid, name, description, price, equipment_type) VALUES
('1', 'Réfrigérateur', 'Réfrigérateur en acier inoxydable', 500.00, '1'),
('2', 'Canapé', "Canapé d'angle en cuir noir", 1000.00, '2'),
('3', 'Téléviseur', 'Téléviseur LED 55 pouces', 800.00, '3');

INSERT INTO HOUSING_EQUIPMENT (number, housing, equipment) VALUES
(1, '1', '1'),
(2, '1', '2'),
(3, '2', '3');

INSERT INTO BED_ROOM (uuid, nbPlaces, price, description, validated, housing) VALUES
('1', 2, 80.00, 'Chambre double avec salle de bain privée', true, '1'),
('2', 1, 50.00, 'Chambre individuelle avec vue sur la ville', false, '1'),
('3', 4, 120.00, 'Suite familiale avec deux chambres', true, '2');

INSERT INTO ACCOUNT_BEDROOM (creation_date, account, bedroom) VALUES
(NOW(), '1', '1'),
(NOW(), '2', '2'),
(NOW(), '3', '3');

INSERT INTO RESERVATION_BEDROOM (start_time, end_time, price, review, review_note, account, bed_room) VALUES
(NOW(), DATE_ADD(NOW(), INTERVAL 7 DAY), 400.00, 'Très bon séjour, chambre confortable', 5, '1', '1'),
(NOW(), DATE_ADD(NOW(), INTERVAL 5 DAY), 250.00, 'Bonne expérience, chambre propre', 4, '2', '2'),
(NOW(), DATE_ADD(NOW(), INTERVAL 10 DAY), 600.00, 'Excellent service, chambre spacieuse', 5, '3', '3');

INSERT INTO RESERVATION_HOUSING (start_time, end_time, price, review, review_note, account, housing) VALUES
(NOW(), DATE_ADD(NOW(), INTERVAL 7 DAY), 1500.00, 'Excellent séjour, logement spacieux', 5, '1', '1'),
(NOW(), DATE_ADD(NOW(), INTERVAL 5 DAY), 1000.00, 'Bonne expérience, logement bien situé', 4, '2', '2'),
(NOW(), DATE_ADD(NOW(), INTERVAL 10 DAY), 2000.00, 'Très satisfait, logement propre et confortable', 5, '3', '3');

INSERT INTO ACCOUNT_HOUSING (creation_date, account, housing) VALUES
(NOW(), '1', '1'),
(NOW(), '2', '2'),
(NOW(), '3', '3');

INSERT INTO MESSAGE (uuid, creation_date, content, note, account, author) VALUES
('1', NOW(), 'Bonjour, je suis intéressé par votre logement.', 4, '1', '2'),
('2', NOW(), 'Bonjour, merci pour votre message.', 5, '2', '1'),
('3', NOW(), 'Pouvons-nous discuter des détails ?', 4, '3', '1');
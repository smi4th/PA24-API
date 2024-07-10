INSERT INTO ACCOUNT_TYPE (uuid, type, private, admin) VALUES
('1', 'Locataire', 'false', 'false'),
('2', 'Loueur', 'false', 'false'),
('3', 'Administateur', 'true', 'true'),
('4', 'Provider', 'true', 'false'),
('5', 'Handyman', 'false', 'false');

INSERT INTO PROVIDER (uuid, name, email, imgPath) VALUES
('1', 'Provider1', 'provider1@provider1.com', NULL),
('2', 'Provider2', 'provider2@provider2.com', NULL),
('3', 'Provider3', 'provider3@provider3.com', NULL);    

INSERT INTO ACCOUNT (uuid, token, username, password, first_name, last_name, email, imgPath, account_type, provider) VALUES
('1', '', 'user1', '$2a$10$8gd1puAaAsl4LBpvAj6uhO4whmoaSKiD69AiZGEamUtDOr6ALZVmG', 'John', 'Doe', 'john.doe@example.com', 'NULL', '1', NULL),
('2', '', 'user2', '$2a$10$rWQKIlPsWD0NIYVjV8xuROlL3ZfZu8nMVq.mIoHuz0fvlYS7aZBda', 'Jane', 'Smith', 'jane.smith@example.com', 'NULL', '2', NULL),
('3', '', 'user3', '$2a$10$Ooru.yReQwZ4v3XoTvQFgu.IgwyuXnaeEqBiVjrdJbWgLlri0juL.', 'Alice', 'Johnson', 'alice.johnson@example.com', 'NULL', '3', NULL),
('4', '', 'user4', '$2a$10$ZPgHF8EJ93a5JLDKCaEmJOhqde6CXepP/NZXhKP1EEbx3kK1BMZgm', 'Bob', 'Brown', 'bob.brown@example.com', 'NULL', '4', '1'),
('5', '', 'user5', '$2a$10$MDVKuspmpvdW/nYaYxpKKe94kYX1bpIj9u40E8GDgOEjBiv/gq5ne', 'Emma', 'Wilson', 'emma.wilson@example.com', 'NULL', '5', '2');

INSERT INTO TAXES(uuid, name, value) VALUES
('1', 'TVA', 20.00),
('2', 'Taxes de séjour', 5.00),
('3', 'Frais de dossier', 10.00);

INSERT INTO SUBSCRIPTION (uuid, name, price, ads, VIP, description, duration, imgPath, taxes) VALUES
('1', 'Subscription1', 10.00, '1', '0', 'Subscription1 description', 30, 'NULL', '1'),
('2', 'Subscription2', 20.00, '0', '0', 'Subscription2 description', 60, 'NULL', '2'),
('3', 'Subscription3', 30.00, '0', '1', 'Subscription3 description', 90, 'NULL', '3');

INSERT INTO ACCOUNT_SUBSCRIPTION (start_date, account, subscription) VALUES
(NOW(), '1', '1'),
(NOW(), '2', '2');

INSERT INTO SERVICES_TYPES (uuid, type, imgPath) VALUES
('1', 'Service1', 'NULL'),
('2', 'Service2', 'NULL'),
('3', 'Service3', 'NULL');

INSERT INTO SERVICES (uuid, price, description, account, service_type, imgPath, duration, token, taxes) VALUES
('1', 10.00, 'Service1 description', '1', '1', 'NULL', '00:30', '123456', '1'),
('2', 20.00, 'Service2 description', '2', '2', 'NULL', '01:00', '654321', '2'),
('3', 30.00, 'Service3 description', '3', '3', 'NULL', '01:30', '987654', '3');

INSERT INTO DISPONIBILITY (uuid, start_date, end_date, account) VALUES
('1', NOW(), DATE_ADD(NOW(), INTERVAL 7 DAY), "1"),
('2', NOW(), DATE_ADD(NOW(), INTERVAL 14 DAY), "2"),
('3', NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), "3");

INSERT INTO HOUSE_TYPE (uuid, type, imgPath) VALUES
('1', 'Maison individuelle', 'NULL'),
('2', 'Appartement', 'NULL'),
('3', 'Maison de ville', 'NULL');

INSERT INTO HOUSING (uuid, surface, price, validated, street_nb, city, zip_code, street, description, house_type, account, imgPath, title, taxes) VALUES
('1', 100.00, 200000.00, true, '123', 'Paris', '75001', 'Rue de Rivoli', 'Belle maison individuelle', '1', '1', 'NULL', 'Maison de charme', '1'),
('2', 75.50, 150000.00, false, '456', 'Lyon', '69001', 'Rue de la République', 'Appartement lumineux', '2', '2', 'NULL', 'Appartement cosy', '1'),
('3', 120.75, 300000.00, true, '789', 'Marseille', '13001', 'Avenue du Prado', 'Maison de ville avec jardin', '3', '1', 'NULL', 'Maison familiale', '1'),
('4', 90.00, 180000.00, true, '321', 'Bordeaux', '33000', 'Rue Sainte-Catherine', 'Charmant appartement au centre-ville', '2', '3', 'NULL', 'Appartement charmant', '1'),
('5', 65.25, 120000.00, false, '654', 'Nice', '06000', 'Promenade des Anglais', 'Studio avec vue sur mer', '1', '4', 'NULL', 'Studio avec vue', '1'),
('6', 110.50, 250000.00, true, '987', 'Lille', '59000', 'Rue Faidherbe', 'Maison spacieuse avec jardin', '3', '2', 'NULL', 'Maison spacieuse', '1'),
('7', 80.00, 160000.00, true, '101', 'Toulouse', '31000', 'Place du Capitole', 'Appartement moderne et lumineux', '2', '5', 'NULL', 'Appartement moderne', '1'),
('8', 95.75, 190000.00, false, '202', 'Nantes', '44000', 'Rue Crébillon', 'Maison rénovée avec terrasse', '3', '1', 'NULL', 'Maison rénovée', '1'),
('9', 70.50, 140000.00, true, '303', 'Strasbourg', '67000', 'Place Kléber', 'Appartement ancien avec charme', '1', '3', 'NULL', 'Appartement ancien', '1'),
('10', 85.00, 170000.00, true, '404', 'Rennes', '35000', 'Rue Saint-Michel', 'Maison de ville bien située', '2', '4', 'NULL', 'Maison de ville', '1'),
('11', 105.25, 220000.00, false, '505', 'Montpellier', '34000', 'Place de la Comédie', 'Grande maison familiale', '3', '2', 'NULL', 'Maison familiale', '1'),
('12', 60.00, 130000.00, true, '606', 'Reims', '51100', 'Rue de Vesle', 'Appartement central et calme', '1', '5', 'NULL', 'Appartement central', '1'),
('13', 115.75, 260000.00, true, '707', 'Dijon', '21000', 'Place de la Libération', 'Maison moderne avec piscine', '3', '1', 'NULL', 'Maison moderne', '1');

INSERT INTO EQUIPMENT_TYPE (uuid, name, imgPath) VALUES
('1', 'Literie', 'NULL'),
('2', 'Mobilier', 'NULL'),
('3', 'Électroménager', 'NULL');

INSERT INTO EQUIPMENT (uuid, name, description, price, equipment_type, housing, imgPath, number, taxes) VALUES
('1', 'Lit double', 'Lit double avec matelas confortable', 100.00, '1', '1', 'NULL', '2', '1'),
('2', 'Canapé', 'Canapé en cuir avec méridienne', 500.00, '2', '2', 'NULL', '1', '2'),
('3', 'Réfrigérateur', 'Réfrigérateur avec congélateur', 800.00, '3', '3', 'NULL', '1', '3');

INSERT INTO BED_ROOM (uuid, nbPlaces, price, description, validated, housing, imgPath, title, taxes) VALUES
('1', 2, 80.00, 'Chambre double avec salle de bain privée', true, '1', 'NULL', 'Chambre parentale', '1'),
('2', 1, 50.00, 'Chambre individuelle avec vue sur la ville', false, '1', 'NULL', 'Chambre avec vue', '1'),
('3', 4, 120.00, 'Suite familiale avec deux chambres', true, '2', 'NULL', 'Suite familiale', '1');

INSERT INTO `BASKET` (uuid, account, paid) VALUES
('1', '1', '0'),
('2', '2', '0'),
('3', '3', '0');

INSERT INTO BASKET_EQUIPMENT (basket, equipment, number) VALUES
('1', '1', 2),
('2', '2', 1),
('3', '3', 3);

INSERT INTO BASKET_BEDROOM (start_time, end_time, basket, bedroom) VALUES
(NOW(), DATE_ADD(NOW(), INTERVAL 7 DAY), '1', '1'),
(NOW(), DATE_ADD(NOW(), INTERVAL 14 DAY), '2', '2'),
(NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), '3', '3');

INSERT INTO BASKET_HOUSING (start_time, end_time, basket, housing) VALUES
(NOW(), DATE_ADD(NOW(), INTERVAL 7 DAY), '1', '1'),
(NOW(), DATE_ADD(NOW(), INTERVAL 14 DAY), '2', '2'),
(NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), '3', '3');

INSERT INTO BASKET_SERVICE (start_time, basket, service) VALUES
(NOW(), '1', '1'),
(NOW(), '2', '2'),
(NOW(), '3', '3');

INSERT INTO REVIEW (uuid, content, note, account, service) VALUES
('1', 'Service rapide et efficace.', 3, '3', '3');
INSERT INTO REVIEW (uuid, content, note, account, housing) VALUES
('2', 'Très bon accueil, je recommande.', 5, '1', '1');
INSERT INTO REVIEW (uuid, content, note, account, bedroom) VALUES
('3', 'Chambre spacieuse et confortable.', 4, '2', '2');

INSERT INTO MESSAGE (uuid, creation_date, content, account, author, imgPath) VALUES
('1', NOW(), 'Bonjour, je suis intéressé par votre logement.', '1', '2', 'NULL'),
('2', NOW(), 'Bonjour, merci pour votre message.', '2', '1', 'NULL'),
('3', NOW(), 'Pouvons-nous discuter des détails ?', '3', '1', 'NULL'),
('4', '2024-06-23 16:20:00', 'Salut, comment ça va ?', '1', '3', 'NULL'),
('5', '2024-06-23 16:21:00', 'Bien, merci. Et toi ?', '2', '3', 'NULL'),
('6', '2024-06-23 16:22:00', 'Je vais bien aussi.', '3', '2', 'NULL'),
('7', '2024-06-23 16:23:00', "J\'ai vu votre annonce, elle m\'intéresse.", '1', '4', 'NULL'),
('8', '2024-06-23 16:24:00', 'Merci pour votre intérêt.', '2', '4', 'NULL'),
('9', '2024-06-23 16:25:00', 'Pouvez-vous me donner plus de détails ?', '3', '4', 'NULL'),
('10', '2024-06-23 16:26:00', 'Bien sûr, quels détails souhaitez-vous ?', '4', '5', 'NULL'),
('11', '2024-06-23 16:27:00', 'Les caractéristiques du logement, par exemple.', '5', '4', 'NULL'),
('12', '2024-06-23 16:28:00', "D\'accord, je vous envoie les informations.", '1', '5', 'NULL'),
('13', '2024-06-23 16:29:00', 'Merci beaucoup.', '5', '1', 'NULL'),
('14', '2024-06-23 16:30:00', 'Vous êtes le bienvenu.', '2', '5', 'NULL'),
('15', '2024-06-23 16:31:00', 'Avez-vous des questions supplémentaires ?', '4', '2', 'NULL'),
('16', '2024-06-23 16:32:00', 'Non, tout est clair. Merci.', '3', '2', 'NULL');

INSERT INTO STATUS (uuid, status) VALUES
('1', 'En attente'),
('2', 'En cours'),
('3', 'Terminé');

INSERT INTO TICKET (uuid, title, description, creation_date, status, account, support) VALUES
('1', 'Problème de connexion', "Je n'arrive pas à me connecter à mon compte.", NOW(), '1', '1', '3'),
('2', 'Problème de paiement', "Je n'arrive pas à payer ma commande.", NOW(), '2', '2', '3'),
('3', 'Problème de réservation', "Je n'arrive pas à réserver un logement.", NOW(), '3', '3', '3');

INSERT INTO TMESSAGE (uuid, content, creation_date, ticket, account) VALUES
('1', 'Avez-vous essayé de réinitialiser votre mot de passe ?', NOW(), '1', '3'),
('2', 'Avez-vous essayé de changer de navigateur ?', NOW(), '2', '3'),
('3', 'Avez-vous essayé de vider votre cache ?', NOW(), '3', '3');

INSERT INTO CHATBOT (uuid, keyword, text) VALUES
(1, 'inscrire', 'Pour créer un compte, cliquez sur ''Se connecter'' en haut à droite et suivez les instructions.'),
(2, 'problème', 'Je suis désolé d''apprendre que vous avez rencontré un problème. Comment puis-je vous aider à le résoudre ?'),
(3, 'appartement', 'Vous cherchez un appartement ? Consultez notre liste d''appartements disponibles en cliquant sur ''voyager''.'),
(4, 'merci', 'Merci pour votre message ! Si vous avez d''autres questions, n''hésitez pas à demander.'),
(5, 'maison', 'Trouvez votre maison idéale en explorant nos options de location sous l''onglet ''voyager''.'),
(6, 'Bonjourww', 'wowowo'),
(7, 'reserver', 'Pour vous renseigner sur les réservations: https://google.com'),
(8, 'connexion', 'Pour se connecter'', cliquez sur ''Se connecter'' en haut à droite et suivez les instructions.'),
(9, 'bonjour', 'Bonjour ! Comment puis-je vous aider aujourd''hui ?'),
(10, 'location', 'Pour toute location, visitez la section ''voyager'' et explorez nos biens disponibles.'),
(11, 'avis', 'Pour lire ou laisser des avis, cliquez sur ''Avis'' en haut de la page.'),
(12, 'au revoir', 'Au revoir ! Passez une excellente journée et revenez nous voir bientôt.'),
(13, 'prestation', 'Pour consulter nos prestations, cliquez sur ''Prestation'' en haut de la page.'),
(14, 'louer', 'Pour louer un bien, cliquez sur ''Louer'' en haut de cette page.'),
(15, 'aide', 'Je suis là pour vous aider. Que puis-je faire pour vous ?'),
(16, 'villa', 'Pour louer une villa, rendez-vous dans la section ''voyager'' et parcourez nos offres exclusives.'),
(17, 'prix', 'Pour connaître les prix et tarifs, rendez-vous dans la section ''voyager''.');
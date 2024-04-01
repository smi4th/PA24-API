DROP DATABASE IF EXISTS `database`;

CREATE DATABASE IF NOT EXISTS `database`;
USE `database`;

CREATE TABLE IF NOT EXISTS `ACCOUNT_TYPE` (
    `id` VARCHAR(40) NOT NULL PRIMARY KEY,
    `type` VARCHAR(45) NOT NULL
);

CREATE TABLE IF NOT EXISTS `ACCOUNT` (
    `id` VARCHAR(40) NOT NULL PRIMARY KEY,
    `username` VARCHAR(45) NOT NULL,
    `password` VARCHAR(45) NOT NULL,
    `first_name` VARCHAR(45) NOT NULL,
    `last_name` VARCHAR(45) NOT NULL,
    `email` VARCHAR(45) NOT NULL,
    `creation_date` DATE DEFAULT NOW(),
    `account_type` VARCHAR(40) NOT NULL,
    FOREIGN KEY (`account_type`) REFERENCES `ACCOUNT_TYPE`(`id`)
);

CREATE TABLE IF NOT EXISTS `SUBSCRIPTION` (
    `id` VARCHAR(40) NOT NULL PRIMARY KEY,
    `name` VARCHAR(45) NOT NULL
);

CREATE TABLE IF NOT EXISTS `ACCOUNT_SUBSCRIPTION` (
    `start_date` DATE DEFAULT NOW(),
    `account` VARCHAR(40) NOT NULL,
    `subscription` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`account`, `subscription`),
    FOREIGN KEY (`account`) REFERENCES `ACCOUNT`(`id`),
    FOREIGN KEY (`subscription`) REFERENCES `SUBSCRIPTION`(`id`)
);

CREATE TABLE IF NOT EXISTS `SERVICES_TYPES` (
    `id` VARCHAR(40) NOT NULL PRIMARY KEY,
    `type` VARCHAR(45) NOT NULL
);

CREATE TABLE IF NOT EXISTS `SERVICES` (
    `id` VARCHAR(40) NOT NULL PRIMARY KEY,
    `price` DECIMAL(10, 2) NOT NULL, -- 10 digits, 2 decimal places
    `service_type` VARCHAR(40) NOT NULL,
    FOREIGN KEY (`service_type`) REFERENCES `SERVICES_TYPES`(`id`)
);

CREATE TABLE IF NOT EXISTS `ACCOUNT_SERVICES` (
    `account` VARCHAR(40) NOT NULL,
    `services` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`account`, `services`),
    FOREIGN KEY (`account`) REFERENCES `ACCOUNT`(`id`),
    FOREIGN KEY (`services`) REFERENCES `SERVICES`(`id`)
);

CREATE TABLE IF NOT EXISTS `consume` (
    `report` TEXT NOT NULL,
    `notice` TEXT NOT NULL,
    `price` DECIMAL(10, 2) NOT NULL,
    `note` INT NOT NULL,
    `services` VARCHAR(40) NOT NULL,
    `account` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`account`, `services`),
    FOREIGN KEY (`account`) REFERENCES `ACCOUNT`(`id`),
    FOREIGN KEY (`services`) REFERENCES `SERVICES`(`id`)
);

CREATE TABLE IF NOT EXISTS `PROVIDER` (
    `id` VARCHAR(40) NOT NULL PRIMARY KEY,
    `name` VARCHAR(45) NOT NULL,
    `email` VARCHAR(45) NOT NULL
);

CREATE TABLE IF NOT EXISTS `PROVIDER_ACCOUNT` (
    `administration_level` INT NOT NULL,
    `provider` VARCHAR(40) NOT NULL,
    `account` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`provider`, `account`),
    FOREIGN KEY (`provider`) REFERENCES `PROVIDER`(`id`),
    FOREIGN KEY (`account`) REFERENCES `ACCOUNT`(`id`)
);

CREATE TABLE IF NOT EXISTS `DISPONIBILITY` (
    `id` VARCHAR(40) NOT NULL PRIMARY KEY,
    `start_date` DATE DEFAULT NOW(),
    `end_date` DATE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS `DISPONIBILITY_ACCOUNT` (
    `disponibility` VARCHAR(40) NOT NULL,
    `account` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`disponibility`, `account`),
    FOREIGN KEY (`disponibility`) REFERENCES `DISPONIBILITY`(`id`),
    FOREIGN KEY (`account`) REFERENCES `ACCOUNT`(`id`)
);

CREATE TABLE IF NOT EXISTS `HOUSE_TYPE` (
    `id` VARCHAR(40) NOT NULL PRIMARY KEY,
    `type` VARCHAR(45) NOT NULL
);

CREATE TABLE IF NOT EXISTS `HOUSING` (
    `id` VARCHAR(40) NOT NULL PRIMARY KEY,
    `surface` DECIMAL(10, 2) NOT NULL,
    `price` DECIMAL(10, 2) NOT NULL,
    `validated` BOOLEAN NOT NULL,
    `street_nb` VARCHAR(45) NOT NULL,
    `city` VARCHAR(45) NOT NULL,
    `zip_code` VARCHAR(45) NOT NULL,
    `street` VARCHAR(45) NOT NULL,
    `description` TEXT NOT NULL,
    `house_type` VARCHAR(40) NOT NULL,
    FOREIGN KEY (`house_type`) REFERENCES `HOUSE_TYPE`(`id`)
);

CREATE TABLE IF NOT EXISTS `EQUIPMENT_TYPE` (
    `id` VARCHAR(40) NOT NULL PRIMARY KEY,
    `name` VARCHAR(45) NOT NULL
);

CREATE TABLE IF NOT EXISTS `EQUIPMENT` (
    `id` VARCHAR(40) NOT NULL PRIMARY KEY,
    `name` VARCHAR(45) NOT NULL,
    `description` TEXT NOT NULL,
    `price` DECIMAL(10, 2) NOT NULL,
    `equipment_type` VARCHAR(40) NOT NULL,
    FOREIGN KEY (`equipment_type`) REFERENCES `EQUIPMENT_TYPE`(`id`)
);

CREATE TABLE IF NOT EXISTS `HOUSING_EQUIPMENT` (
    `number` INT NOT NULL,
    `housing` VARCHAR(40) NOT NULL,
    `equipment` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`housing`, `equipment`),
    FOREIGN KEY (`housing`) REFERENCES `HOUSING`(`id`),
    FOREIGN KEY (`equipment`) REFERENCES `EQUIPMENT`(`id`)
);

CREATE TABLE IF NOT EXISTS `BED_ROOM` (
    `id` VARCHAR(40) NOT NULL PRIMARY KEY,
    `nbPlaces` INT NOT NULL,
    `price` DECIMAL(10, 2) NOT NULL,
    `description` TEXT NOT NULL,
    `validated` BOOLEAN NOT NULL,
    `housing` VARCHAR(40) NOT NULL,
    FOREIGN KEY (`housing`) REFERENCES `HOUSING`(`id`)
);

CREATE TABLE IF NOT EXISTS `ACCOUNT_BEDROOM` (
    `creation_date` DATE DEFAULT NOW(),
    `account` VARCHAR(40) NOT NULL,
    `bedroom` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`account`, `bedroom`),
    FOREIGN KEY (`account`) REFERENCES `ACCOUNT`(`id`),
    FOREIGN KEY (`bedroom`) REFERENCES `BED_ROOM`(`id`)
);

CREATE TABLE IF NOT EXISTS `RESERVATION_BEDROOM` (
    `start_time` DATE DEFAULT NOW(),
    `end_time` DATE DEFAULT NOW(),
    `price` DECIMAL(10, 2) NOT NULL,
    `review` TEXT NOT NULL,
    `review_note` INT NOT NULL,
    `account` VARCHAR(40) NOT NULL,
    `bed_room` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`account`, `bed_room`),
    FOREIGN KEY (`account`) REFERENCES `ACCOUNT`(`id`),
    FOREIGN KEY (`bed_room`) REFERENCES `BED_ROOM`(`id`)
);

CREATE TABLE IF NOT EXISTS `RESERVATION_HOUSING` (
    `start_time` DATE DEFAULT NOW(),
    `end_time` DATE DEFAULT NOW(),
    `price` DECIMAL(10, 2) NOT NULL,
    `review` TEXT NOT NULL,
    `review_note` INT NOT NULL,
    `account` VARCHAR(40) NOT NULL,
    `housing` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`account`, `housing`),
    FOREIGN KEY (`account`) REFERENCES `ACCOUNT`(`id`),
    FOREIGN KEY (`housing`) REFERENCES `HOUSING`(`id`)
);

CREATE TABLE IF NOT EXISTS `ACCOUNT_HOUSING` (
    `creation_date` DATE DEFAULT NOW(),
    `account` VARCHAR(40) NOT NULL,
    `housing` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`account`, `housing`),
    FOREIGN KEY (`account`) REFERENCES `ACCOUNT`(`id`),
    FOREIGN KEY (`housing`) REFERENCES `HOUSING`(`id`)
);

CREATE TABLE IF NOT EXISTS `MESSAGE` (
    `creation_date` DATE DEFAULT NOW(),
    `content` TEXT NOT NULL,
    `note` INT NOT NULL,
    `account` VARCHAR(40) NOT NULL,
    `author` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`account`, `author`, `creation_date`),
    FOREIGN KEY (`account`) REFERENCES `ACCOUNT`(`id`),
    FOREIGN KEY (`author`) REFERENCES `ACCOUNT`(`id`)
);
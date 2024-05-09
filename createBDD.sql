DROP DATABASE IF EXISTS `database`;

CREATE DATABASE IF NOT EXISTS `database`;
USE `database`;

CREATE TABLE IF NOT EXISTS `ACCOUNT_TYPE` (
    `uuid` VARCHAR(40) NOT NULL PRIMARY KEY,
    `type` VARCHAR(45) NOT NULL, -- UNIQUE
    `private` CHAR(5) NOT NULL DEFAULT 'false',
    `admin` CHAR(5) NOT NULL DEFAULT 'false'
);

CREATE TABLE IF NOT EXISTS `PROVIDER` (
    `uuid` VARCHAR(40) NOT NULL PRIMARY KEY,
    `name` VARCHAR(45) NOT NULL, -- UNIQUE
    `email` VARCHAR(45) NOT NULL, -- UNIQUE
    `imgPath` VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS `ACCOUNT` (
    `uuid` VARCHAR(40) NOT NULL PRIMARY KEY,
    `token` VARCHAR(64),
    `username` VARCHAR(45) NOT NULL, -- UNIQUE
    `password` VARCHAR(60) NOT NULL,
    `first_name` VARCHAR(45) NOT NULL,
    `last_name` VARCHAR(45) NOT NULL,
    `email` VARCHAR(45) NOT NULL, -- UNIQUE
    `creation_date` DATE DEFAULT NOW(), -- AUTO GEN
    `imgPath` VARCHAR(255),
    `account_type` VARCHAR(40) NOT NULL,
    `provider` VARCHAR(40),
    FOREIGN KEY (`account_type`) REFERENCES `ACCOUNT_TYPE`(`uuid`),
    FOREIGN KEY (`provider`) REFERENCES `PROVIDER`(`uuid`)
);

CREATE TABLE IF NOT EXISTS `SUBSCRIPTION` (
    `uuid` VARCHAR(40) NOT NULL PRIMARY KEY,
    `name` VARCHAR(45) NOT NULL, -- UNIQUE
    `price` DECIMAL(10, 2) NOT NULL,
    `ads` BOOLEAN NOT NULL,
    `VIP` BOOLEAN NOT NULL,
    `description` TEXT NOT NULL,
    `duration` INT NOT NULL,
    `imgPath` VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS `ACCOUNT_SUBSCRIPTION` (
    `start_date` DATE DEFAULT NOW(), -- AUTO GEN
    `account` VARCHAR(40) NOT NULL,
    `subscription` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`account`, `subscription`),
    FOREIGN KEY (`account`) REFERENCES `ACCOUNT`(`uuid`),
    FOREIGN KEY (`subscription`) REFERENCES `SUBSCRIPTION`(`uuid`)
);

CREATE TABLE IF NOT EXISTS `SERVICES_TYPES` (
    `uuid` VARCHAR(40) NOT NULL PRIMARY KEY,
    `type` VARCHAR(45) NOT NULL, -- UNIQUE
    `imgPath` VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS `SERVICES` (
    `uuid` VARCHAR(40) NOT NULL PRIMARY KEY,
    `price` DECIMAL(10, 2) NOT NULL,
    `description` TEXT NOT NULL,
    `imgPath` VARCHAR(255),
    `duration` TIME NOT NULL,
    `account` VARCHAR(40) NOT NULL,
    `service_type` VARCHAR(40) NOT NULL,
    FOREIGN KEY (`account`) REFERENCES `ACCOUNT`(`uuid`),
    FOREIGN KEY (`service_type`) REFERENCES `SERVICES_TYPES`(`uuid`)
);

CREATE TABLE IF NOT EXISTS `DISPONIBILITY` (
    `uuid` VARCHAR(40) NOT NULL PRIMARY KEY,
    `start_date` DATE DEFAULT NOW(),
    `end_date` DATE DEFAULT NOW(),
    `account` VARCHAR(40) NOT NULL,
    FOREIGN KEY (`account`) REFERENCES `ACCOUNT`(`uuid`)
);

CREATE TABLE IF NOT EXISTS `HOUSE_TYPE` (
    `uuid` VARCHAR(40) NOT NULL PRIMARY KEY,
    `type` VARCHAR(45) NOT NULL, -- UNIQUE
    `imgPath` VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS `HOUSING` (
    `uuid` VARCHAR(40) NOT NULL PRIMARY KEY,
    `surface` DECIMAL(10, 2) NOT NULL,
    `price` DECIMAL(10, 2) NOT NULL,
    `validated` BOOLEAN NOT NULL,
    `street_nb` VARCHAR(45) NOT NULL,
    `city` VARCHAR(45) NOT NULL,
    `zip_code` VARCHAR(45) NOT NULL,
    `street` VARCHAR(45) NOT NULL,
    `description` TEXT NOT NULL,
    `imgPath` VARCHAR(255),
    `house_type` VARCHAR(40) NOT NULL,
    `account` VARCHAR(40) NOT NULL,
    FOREIGN KEY (`house_type`) REFERENCES `HOUSE_TYPE`(`uuid`),
    FOREIGN KEY (`account`) REFERENCES `ACCOUNT`(`uuid`)
);

CREATE TABLE IF NOT EXISTS `EQUIPMENT_TYPE` (
    `uuid` VARCHAR(40) NOT NULL PRIMARY KEY,
    `name` VARCHAR(45) NOT NULL, -- UNIQUE
    `imgPath` VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS `EQUIPMENT` (
    `uuid` VARCHAR(40) NOT NULL PRIMARY KEY,
    `name` VARCHAR(45) NOT NULL,
    `description` TEXT NOT NULL,
    `price` DECIMAL(10, 2) NOT NULL,
    `number` INT NOT NULL,
    `imgPath` VARCHAR(255),
    `equipment_type` VARCHAR(40) NOT NULL,
    `housing` VARCHAR(40) NOT NULL,
    FOREIGN KEY (`equipment_type`) REFERENCES `EQUIPMENT_TYPE`(`uuid`),
    FOREIGN KEY (`housing`) REFERENCES `HOUSING`(`uuid`)
);

CREATE TABLE IF NOT EXISTS `BED_ROOM` (
    `uuid` VARCHAR(40) NOT NULL PRIMARY KEY,
    `nbPlaces` INT NOT NULL,
    `price` DECIMAL(10, 2) NOT NULL,
    `description` TEXT NOT NULL,
    `validated` BOOLEAN NOT NULL,
    `imgPath` VARCHAR(255),
    `housing` VARCHAR(40) NOT NULL,
    FOREIGN KEY (`housing`) REFERENCES `HOUSING`(`uuid`)
);

CREATE TABLE IF NOT EXISTS `BASKET` (
    `uuid` VARCHAR(40) NOT NULL PRIMARY KEY,
    `account` VARCHAR(40) NOT NULL,
    `paid` BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS `BASKET_EQUIPMENT` (
    `basket` VARCHAR(40) NOT NULL,
    `equipment` VARCHAR(40) NOT NULL,
    `number` INT NOT NULL,
    PRIMARY KEY (`basket`, `equipment`),
    FOREIGN KEY (`basket`) REFERENCES `BASKET`(`uuid`),
    FOREIGN KEY (`equipment`) REFERENCES `EQUIPMENT`(`uuid`)
);

CREATE TABLE IF NOT EXISTS `BASKET_BEDROOM` (
    `start_time` DATE DEFAULT NOW(),
    `end_time` DATE DEFAULT NOW(),
    `basket` VARCHAR(40) NOT NULL,
    `bedroom` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`basket`, `bedroom`),
    FOREIGN KEY (`basket`) REFERENCES `BASKET`(`uuid`),
    FOREIGN KEY (`bedroom`) REFERENCES `BED_ROOM`(`uuid`)
);

CREATE TABLE IF NOT EXISTS `BASKET_HOUSING` (
    `start_time` DATE DEFAULT NOW(),
    `end_time` DATE DEFAULT NOW(),
    `basket` VARCHAR(40) NOT NULL,
    `housing` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`basket`, `housing`),
    FOREIGN KEY (`basket`) REFERENCES `BASKET`(`uuid`),
    FOREIGN KEY (`housing`) REFERENCES `HOUSING`(`uuid`)
);

CREATE TABLE IF NOT EXISTS `BASKET_SERVICE` (
    `start_time` DATETIME DEFAULT NOW(),
    `basket` VARCHAR(40) NOT NULL,
    `service` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`basket`, `service`),
    FOREIGN KEY (`basket`) REFERENCES `BASKET`(`uuid`),
    FOREIGN KEY (`service`) REFERENCES `SERVICES`(`uuid`)
);

CREATE TABLE IF NOT EXISTS `REVIEW` (
    `uuid` VARCHAR(40) NOT NULL PRIMARY KEY,
    `content` TEXT NOT NULL,
    `note` INT NOT NULL,
    `account` VARCHAR(40) NOT NULL,
    `service` VARCHAR(40) NULL,
    `housing` VARCHAR(40) NULL,
    `bedroom` VARCHAR(40) NULL,
    FOREIGN KEY (`account`) REFERENCES `ACCOUNT`(`uuid`),
    FOREIGN KEY (`service`) REFERENCES `SERVICES`(`uuid`),
    FOREIGN KEY (`housing`) REFERENCES `HOUSING`(`uuid`),
    FOREIGN KEY (`bedroom`) REFERENCES `BED_ROOM`(`uuid`)
);

CREATE TABLE IF NOT EXISTS `MESSAGE` (
    `uuid` VARCHAR(40) NOT NULL,
    `creation_date` DATE DEFAULT NOW(), -- AUTO GEN
    `content` TEXT NOT NULL,
    `imgPath` VARCHAR(255),
    `account` VARCHAR(40) NOT NULL,
    `author` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`account`, `author`, `uuid`),
    FOREIGN KEY (`account`) REFERENCES `ACCOUNT`(`uuid`),
    FOREIGN KEY (`author`) REFERENCES `ACCOUNT`(`uuid`)
);
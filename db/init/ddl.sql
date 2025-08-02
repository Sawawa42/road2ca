CREATE DATABASE IF NOT EXISTS `road2ca` DEFAULT CHARACTER SET utf8mb4;
USE `road2ca`;

SET CHARSET utf8mb4;

-- 設定情報を格納するテーブル
CREATE TABLE IF NOT EXISTS `road2ca`.`settings` (
    `id` BINARY(16) NOT NULL,
    `name` VARCHAR(128) NOT NULL,
    `gachaCoinConsumption` INT NOT NULL DEFAULT 0,
    `drawGachaMaxTimes` INT NOT NULL DEFAULT 0,
    `getRankingLimit` INT NOT NULL DEFAULT 0,
    `rewardCoin` INT NOT NULL DEFAULT 0,
    `rarity3Ratio` INT NOT NULL DEFAULT 0,
    `rarity2Ratio` INT NOT NULL DEFAULT 0,
    `rarity1Ratio` INT NOT NULL DEFAULT 0,
    `createdAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updatedAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);

-- ユーザー情報を格納するテーブル
CREATE TABLE IF NOT EXISTS `road2ca`.`users` (
    `id` BINARY(16) NOT NULL,
    `name` VARCHAR(128) NOT NULL,
    `highscore` INT NOT NULL DEFAULT 0,
    `coin` INT NOT NULL DEFAULT 0,
    `token` VARCHAR(128) NOT NULL UNIQUE,
    `createdAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updatedAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);

-- アイテム情報を格納するテーブル
CREATE TABLE IF NOT EXISTS `road2ca`.`items` (
    `id` BINARY(16) NOT NULL,
    `name` VARCHAR(128) NOT NULL,
    `rarity` TINYINT NOT NULL DEFAULT 0,
    `weight` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`)
);

-- ユーザの持つアイテム情報を格納する中間テーブル
CREATE TABLE IF NOT EXISTS `road2ca`.`collections` (
    `id` BINARY(16) NOT NULL,
    `userId` BINARY(16) NOT NULL,
    `itemId` BINARY(16) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`userId`) REFERENCES `road2ca`.`users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (`itemId`) REFERENCES `road2ca`.`items`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

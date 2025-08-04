use `road2ca`;

SET CHARSET utf8mb4;

INSERT INTO `settings` (
    `name`,
    `gachaCoinConsumption`,
    `drawGachaMaxTimes`,
    `getRankingLimit`,
    `rewardCoin`,
    `rarity3Ratio`,
    `rarity2Ratio`,
    `rarity1Ratio`
) VALUES (
    'default',
    100,
    10,
    10,
    1000,
    0.6,
    5.1,
    94.3
);

INSERT INTO `items` (`name`, `rarity`) VALUES
('木の棒', 1),
('布の服', 1),
('紙の兜', 1),
('鉄の剣', 2),
('鉄の鎧', 2),
('鉄の兜', 2),
('nop', 3),
('sho', 3),
('15℃', 3);

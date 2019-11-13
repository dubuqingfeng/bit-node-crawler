CREATE TABLE `peers` (
  `address` varchar(255) NOT NULL DEFAULT '',
  `height` bigint(20) NOT NULL DEFAULT '0',
  `peers` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `user_agent` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `type` varchar(255) NOT NULL DEFAULT '',
  `prev_hash` varchar(255) NOT NULL DEFAULT '',
  `coin_type` varchar(50) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL DEFAULT '',
  `timestamp` datetime DEFAULT NULL,
  `notified_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`address`),
  KEY `height` (`height`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
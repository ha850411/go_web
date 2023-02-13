CREATE TABLE `menu` (
  `id` int(11) NOT NULL,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `icon` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `path` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `active` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sort` int(11) NOT NULL,
  `updateTime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `menu` VALUES 
(1,'存貨管理','fa-solid fa-cart-shopping','/product','product',1,'2023-02-10 11:07:15'),
(2,'存貨分析','fa-solid fa-chart-simple','/analysis','analysis',2,'2023-02-10 11:07:22'),
(3,'訂單管理','fa-solid fa-list','/order','order',3,'2023-02-10 11:17:56'),
(4,'設定','fa-solid fa-gear','/setting','setting',4,'2023-02-10 11:07:28');

CREATE TABLE `orders` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `contact` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `remark` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `updateTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `products` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `amount` int(11) NOT NULL,
  `amountNotice` int(11) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '0:刪除 1:存在',
  `updateTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `products_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pid` int(11) NOT NULL,
  `amount` int(11) NOT NULL,
  `updateDate` date NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `pid_2` (`pid`,`updateDate`),
  KEY `pid` (`pid`),
  KEY `pid_3` (`pid`,`updateDate`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `linebot_token` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'linebot access_token',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `users` VALUES 
(1,'admin','admin',NULL),
(2,'eason','eason',NULL),
(3,'niigataiwan','niigataiwan',NULL);

CREATE TABLE `orders_detail` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `order_id` int(11) NOT NULL,
  `pid` int(11) NOT NULL,
  `amount` int(11) NOT NULL,
  `updateTime` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `order_id` (`order_id`),
  KEY `pid` (`pid`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
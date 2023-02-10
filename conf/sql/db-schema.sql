SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";

CREATE TABLE `products` (
  `id` int NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `amount` int NOT NULL,
  `amountNotice` int NOT NULL DEFAULT '0',
  `status` int NOT NULL DEFAULT '1' COMMENT '0:刪除 1:存在',
  `updateTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `products`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `products`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
--------------------------------------------------------------------------------------------------

CREATE TABLE `products_log` (
  `id` int NOT NULL,
  `pid` int NOT NULL,
  `amount` int NOT NULL,
  `updateDate` date NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `products_log`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `pid_2` (`pid`,`updateDate`),
  ADD KEY `pid` (`pid`),
  ADD KEY `pid_3` (`pid`,`updateDate`);

ALTER TABLE `products_log`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
--------------------------------------------------------------------------------------------------

CREATE TABLE `users` (
  `id` int NOT NULL,
  `username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `users`
  ADD PRIMARY KEY (`id`);
  
ALTER TABLE `users`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--------------------------------------------------------------------------------------------------
CREATE TABLE `menu` (
  `id` int NOT NULL,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `icon` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `path` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `active` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sort` int NOT NULL,
  `updateTime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `menu`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `menu`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;
--------------------------------------------------------------------------------------------------
CREATE TABLE `orders` (
  `id` int NOT NULL,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `contact` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pid` int NOT NULL,
  `amount` int NOT NULL,
  `remark` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `updateTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `orders`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `orders`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
--------------------------------------------------------------------------------------------------

COMMIT;
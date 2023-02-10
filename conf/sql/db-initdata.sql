INSERT INTO `users` (`id`, `username`, `password`, `update_time`) VALUES
(1, 'admin', 'admin', NULL),
(2, 'eason', 'eason', NULL),
(3, 'niigataiwan', 'niigataiwan', NULL);

INSERT INTO `menu` (`id`, `name`, `icon`, `path`, `active`, `sort`, `updateTime`) VALUES
(1, '存貨管理', 'fa-solid fa-cart-shopping', '/product', 'product', 1, '2023-02-10 11:07:15'),
(2, '存貨分析', 'fa-solid fa-chart-simple', '/analysis', 'analysis', 2, '2023-02-10 11:07:22'),
(3, '訂單管理', 'fa-solid fa-list', '/order', 'order', 3, '2023-02-10 11:17:56'),
(4, '設定', 'fa-solid fa-gear', '/setting', 'setting', 4, '2023-02-10 11:07:28');
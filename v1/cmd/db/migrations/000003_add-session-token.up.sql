ALTER TABLE `users`
ADD `refresh_token` varchar(40) NULL,
ADD `token_expire` int(64) NULL;
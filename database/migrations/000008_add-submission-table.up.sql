CREATE TABLE IF NOT EXISTS `submissions` (
  `id` int(64) NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `user_id` int(10) NOT NULL,
  `proposal_id` int(64) NOT NULL,
  `file_name` varchar(100) NOT NULL,
  `created_at` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  `updated_at` datetime NULL AFTER `created_at`,
  `deleted_at` datetime NULL AFTER `updated_at`,
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  FOREIGN KEY (`proposal_id`) REFERENCES `proposals` (`id`)
);
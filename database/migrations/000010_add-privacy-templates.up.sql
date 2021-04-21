CREATE TABLE IF NOT EXISTS`privacy_templates` (
  `created_at` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  `updated_at` datetime NULL AFTER `created_at`,
  `deleted_at` datetime NULL AFTER `updated_at`,
  `id` int(64) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `content` longtext NOT NULL,
  `name` varchar(100) COLLATE 'utf8mb4_general_ci' NOT NULL;
);
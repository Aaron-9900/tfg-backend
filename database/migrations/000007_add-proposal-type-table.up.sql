CREATE TABLE IF NOT EXISTS `proposal_types` (
  `id` int(64) NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `value` varchar(32) NOT NULL
);
ALTER TABLE `proposal_types`
ADD `created_at` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
ADD `updated_at` datetime NOT NULL AFTER `created_at`,
ADD `deleted_at` datetime NOT NULL AFTER `updated_at`;
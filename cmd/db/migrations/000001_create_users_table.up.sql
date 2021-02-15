CREATE TABLE IF NOT EXISTS `users` (
	`id` int(10) NOT NULL auto_increment,
	`name` varchar(255),
	`email` varchar(255),
	`password` varchar(255),
	PRIMARY KEY( `id` )
);
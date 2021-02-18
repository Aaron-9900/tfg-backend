CREATE TABLE IF NOT EXISTS `proposals` (
	id int(64) auto_increment NOT NULL,
	user_id int(64) NOT NULL,
	name varchar(100) NULL,
	description TEXT NULL,
	`limit` int(64) NULL,
	CONSTRAINT proposal_PK PRIMARY KEY (id),
	CONSTRAINT proposal_FK FOREIGN KEY (user_id) REFERENCES `users`(id) ON DELETE CASCADE
)
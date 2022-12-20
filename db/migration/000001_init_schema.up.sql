CREATE TABLE `credentials` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `username` varchar(64) NOT NULL,
  `password` varchar(64) NULL,
  `email` varchar(64) NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `token` (
    `id` bigint PRIMARY KEY AUTO_INCREMENT,
	`credentials_id` bigint NOT NULL,
	`token` varchar(256) NULL,
	`created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE `token` ADD FOREIGN KEY (`credentials_id`) REFERENCES `credentials` (`id`);

CREATE INDEX credentials_username_password_idx ON `credentials` (`username`, `password`);

CREATE INDEX credentials_email_password_idx ON `credentials` (`email`, `password`);
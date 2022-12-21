CREATE TABLE `role` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `enabled` bit NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime
);

CREATE TABLE `action` (
    `id` bigint PRIMARY KEY AUTO_INCREMENT,
	`action` varchar(64) NOT NULL,
    `entity` varchar(64) NOT NULL,
    `enabled` bit NOT NULL DEFAULT 1,
	`created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime
);

CREATE TABLE `roles_actions` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `role_id` bigint NOT NULL,
  `action_id` bigint NULL,
  constraint role_id_fk_1
    foreign key (role_id) references `role` (id),
  constraint action_id_fk_1
    foreign key (action_id) references `action` (id)
);

CREATE TABLE `roles_credentials` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `role_id` bigint NOT NULL,
  `credentials_id` bigint NULL,
  constraint role_id_fk_2
    foreign key (role_id) references `role` (id),
  constraint credentials_id_fk_1
    foreign key (credentials_id) references `credentials` (id)
);
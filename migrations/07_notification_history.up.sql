CREATE TABLE IF NOT EXISTS `notification_history` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `recipient` VARCHAR(255) NOT NULL,
    `template_id` INT NOT NULL,
    `status` VARCHAR(10) NOT NULL,
    `error_message` TEXT,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_notification_history_template`
        FOREIGN KEY (`template_id`)
        REFERENCES `notification_templates` (`id`)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

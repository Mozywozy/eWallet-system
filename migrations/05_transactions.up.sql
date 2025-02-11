CREATE TABLE IF NOT EXISTS `transactions` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `user_id` INT NOT NULL,
    `amount` DECIMAL(15,2) NOT NULL DEFAULT 0,
    `transaction_type` ENUM('TOPUP','PURCHASE','REFUND') NOT NULL,
    `transaction_status` ENUM('PENDING','SUCCESS','FAILED','REVERSED') NOT NULL DEFAULT 'PENDING',
    `reference` VARCHAR(255) NOT NULL,
    `description` VARCHAR(255) NOT NULL,
    `additional_info` JSON,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
        CONSTRAINT `fk_transactions_user_id` 
        FOREIGN KEY (`user_id`) 
        REFERENCES `users` (`id`) 
        ON UPDATE CASCADE 
        ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

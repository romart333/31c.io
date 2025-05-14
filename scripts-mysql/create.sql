DROP TABLE IF EXISTS account;
CREATE TABLE account (
    id BIGINT NOT NULL AUTO_INCREMENT,
    address VARCHAR(64) NOT NULL,
    balance DECIMAL(64, 8) NOT NULL DEFAULT 0,
    status ENUM('On', 'Off') NOT NULL,
    created_at INT NOT NULL,
    updated_at INT NOT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX account_address_unique_idx (address),
    INDEX account_status_idx (status),
    INDEX account_updated_idx (updated_at)
);

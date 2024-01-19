START TRANSACTION;

CREATE TABLE accounts
(
    id         VARCHAR(32) PRIMARY KEY,
    name       VARCHAR(64)                         NOT NULL,
    status     VARCHAR(32)                         NOT NULL,
    phone_no   VARCHAR(10)                         NOT NULL,
    gender     VARCHAR(10)                         NOT NULL,
    address    VARCHAR(128)                         NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE student_parents
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    parent_id  VARCHAR(32)                         NOT NULL,
    student_id VARCHAR(32)                         NOT NULL,
    status     VARCHAR(32)                         NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE registers
(
    id            VARCHAR(32) PRIMARY KEY,
    parent_id     VARCHAR(32)                         NOT NULL,
    student_id    VARCHAR(32)                         NOT NULL,
    status        VARCHAR(32)                         NOT NULL,
    register_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE transactions
(
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    register_id VARCHAR(32)                         NOT NULL,
    action_type VARCHAR(64)                         NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

COMMIT;
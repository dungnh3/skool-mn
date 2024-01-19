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

CREATE UNIQUE INDEX student_parents_uindex on student_parents (parent_id, student_id);

CREATE INDEX student_parents_parent_id_index on student_parents (parent_id);

CREATE INDEX student_parents_student_id_index on student_parents (student_id);

CREATE INDEX student_parents_status_index on student_parents (status);

ALTER TABLE student_parents
    ADD CONSTRAINT student_id_fkey FOREIGN KEY (student_id) REFERENCES accounts (id);

ALTER TABLE student_parents
    ADD CONSTRAINT parent_id_fkey FOREIGN KEY (parent_id) REFERENCES accounts (id);

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

CREATE INDEX registers_parent_id_index on registers (parent_id);

CREATE INDEX registers_student_id_index on registers (student_id);

CREATE INDEX registers_status_index on registers (status);

CREATE TABLE transactions
(
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    register_id VARCHAR(32)                         NOT NULL,
    action_type VARCHAR(64)                         NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX transactions_register_id_index on transactions (register_id);

CREATE INDEX transactions_action_type_index on transactions (action_type);

COMMIT;
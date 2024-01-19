# Migration

```mysql
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

INSERT INTO skool_mn.accounts (id, name, status, phone_no, gender, address, created_at, updated_at) VALUES ('17E8ZIN', 'Nguyen Van A1', 'active', '0933123457', 'male', '{"city": "Hồ Chí Minh", "address": "285 CMT8", "district": "Q.10"}', '2024-01-18 22:30:55', '2024-01-18 22:30:55');
INSERT INTO skool_mn.accounts (id, name, status, phone_no, gender, address, created_at, updated_at) VALUES ('1BX7NAJ', 'Nguyen Van A', 'active', '0933123456', 'male', '{"city": "Hồ Chí Minh", "address": "285 CMT8", "district": "Q.10"}', '2024-01-18 22:30:55', '2024-01-18 22:30:55');
INSERT INTO skool_mn.accounts (id, name, status, phone_no, gender, address, created_at, updated_at) VALUES ('1EYTO4N', 'Le Thi B', 'active', '0933123459', 'female', '{"city": "Hồ Chí Minh", "address": "285 CMT8", "district": "Q.10"}', '2024-01-18 22:30:55', '2024-01-18 15:48:01');
INSERT INTO skool_mn.accounts (id, name, status, phone_no, gender, address, created_at, updated_at) VALUES ('1YZB7F4', 'Nguyen Thi A2', 'active', '0933123458', 'female', '{"city": "Hồ Chí Minh", "address": "285 CMT8", "district": "Q.10"}', '2024-01-18 22:30:55', '2024-01-18 15:48:01');


COMMIT;
```

CREATE SCHEMA IF NOT EXISTS ssego DEFAULT CHARACTER SET utf8mb4;

USE ssego;

CREATE TABLE IF NOT EXISTS documents (
  PRIMARY KEY (document_id),
  document_id    INT UNSIGNED AUTO_INCREMENT NOT NULL,
  document_title TEXT                        NOT NULL,
  updated_at     DATETIME default current_timestamp on update current_timestamp,
  created_at     DATETIME default current_timestamp
) ENGINE=InnoDB DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_bin;

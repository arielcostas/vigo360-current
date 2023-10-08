USE vigo360;

CREATE TABLE log(
    rid CHAR(15) NOT NULL,
    sid CHAR(15) NOT NULL,
    time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ip VARCHAR(40) NOT NULL,
    url VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    time_taken_ms INT NOT NULL,
    user_agent VARCHAR(255) NOT NULL,
    PRIMARY KEY(rid),
    INDEX (time_taken_ms),
    INDEX (url)
);


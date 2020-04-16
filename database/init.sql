CREATE DATABASE IF NOT EXISTS ElasticJury DEFAULT CHARACTER SET utf8;
USE ElasticJury;
CREATE TABLE IF NOT EXISTS Cases -- 案件数据库
(
    `id`      INT UNSIGNED  NOT NULL AUTO_INCREMENT,
    `date`    DATE          NULL,     -- TODO 案发时间？一审时间？
    `court`   CHAR(32)      NULL,     -- 法院
    `parties` VARCHAR(512)  NULL,     -- 当事人信息：原告(plaintiff) 被告(defendant) 法定代表人...
    `judge`   CHAR(32)      NULL,     -- 法官
    `law`     VARCHAR(1024) NULL,     -- 法律，考虑到一对多关系，倒排索引需单独建表
    `tag`     CHAR(128)     NULL,     -- 标签，考虑到一对多关系，倒排索引需单独建表
    `detail`  LONGTEXT      NOT NULL, -- 案情 TODO discussion on choice of data type
    PRIMARY KEY (`id` ASC),
    INDEX `date_idx` (`date` ASC),    -- TODO discussion
    INDEX `court_idx` (`court` ASC),
    INDEX `judge_idx` (`judge` ASC)
);
CREATE TABLE IF NOT EXISTS Words -- 词数据库, word -> wordId
(
    `id`   INT UNSIGNED NOT NULL AUTO_INCREMENT, -- wordId
    `word` CHAR(32)     NOT NULL,                -- 去停用词，词根化
    PRIMARY KEY (`id` ASC),
    INDEX `word_idx` (`word` ASC)                -- TODO discussion
);
CREATE TABLE IF NOT EXISTS WordCase -- 词倒排索引数据库（一对多）, wordId -> [](caseId, weight)
(
    `id`     INT UNSIGNED   NOT NULL AUTO_INCREMENT, -- 索引条目的id
    `wordId` INT UNSIGNED   NOT NULL,
    `caseId` INT UNSIGNED   NOT NULL,
    `weight` FLOAT UNSIGNED NOT NULL,
    PRIMARY KEY (`id` ASC),
    INDEX `wordId_idx` (`wordId` ASC)                -- TODO discussion 有套娃嫌疑，但似乎必须有这个字段？
);
CREATE TABLE IF NOT EXISTS Laws -- Law -> LawId
(
    `id`  INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `law` CHAR(64)     NOT NULL,
    PRIMARY KEY (`id` ASC),
    INDEX `law_idx` (`law` ASC)
);
CREATE TABLE IF NOT EXISTS LawCase
(
    `id`     INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `lawId`  INT UNSIGNED NOT NULL,
    `caseId` INT UNSIGNED NOT NULL,
    PRIMARY KEY (`id` ASC),
    INDEX `lawId_idx` (`lawId` ASC)
);
CREATE TABLE IF NOT EXISTS Tags -- Tag -> TagId
(
    `id`  INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `tag` CHAR(64)     NOT NULL,
    PRIMARY KEY (`id` ASC),
    INDEX `tag_idx` (`tag` ASC)
);
CREATE TABLE IF NOT EXISTS TagCase
(
    `id`     INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `tagId`  INT UNSIGNED NOT NULL,
    `caseId` INT UNSIGNED NOT NULL,
    PRIMARY KEY (`id` ASC),
    INDEX `caseId_idx` (`tagId` ASC)
);

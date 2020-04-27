CREATE TABLE IF NOT EXISTS Cases # 案件数据库
(
    `id`     INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `title`  TINYTEXT     NOT NULL, # 标题
    `judge`  TEXT         NULL,     # 法官，考虑到一对多关系，倒排索引需单独建表
    `law`    TEXT         NULL,     # 法律，考虑到一对多关系，倒排索引需单独建表
    `tag`    TEXT         NULL,     # 标签/关键词，考虑到一对多关系，倒排索引需单独建表
    `link`   TEXT         NULL,     # 原网页链接
    `detail` LONGTEXT     NOT NULL, # 案情（xml 中的全文部分）
    `tree`	 LONGTEXT	  NOT NULL, # xml 的树形结构
    PRIMARY KEY (`id` ASC)
) CHAR SET utf8;
CREATE TABLE IF NOT EXISTS WordIndex # 词倒排索引数据库（一对多）, word -> [](case, weight)
(
    `word`   VARCHAR(64)  NOT NULL,
    `caseId` INT UNSIGNED NOT NULL,
    `weight` FLOAT        NOT NULL,
    FOREIGN KEY (`caseId`) REFERENCES Cases (`id`)
) CHAR SET utf8;
CREATE TABLE IF NOT EXISTS JudgeIndex
(
    `judge`  VARCHAR(64)  NOT NULL,
    `caseId` INT UNSIGNED NOT NULL,
    `weight` FLOAT        NOT NULL,
    FOREIGN KEY (`caseId`) REFERENCES Cases (`id`)
) CHAR SET utf8;
CREATE TABLE IF NOT EXISTS LawIndex
(
    `law`    TEXT         NOT NULL,
    `caseId` INT UNSIGNED NOT NULL,
    `weight` FLOAT        NOT NULL,
    FOREIGN KEY (`caseId`) REFERENCES Cases (`id`)
) CHAR SET utf8;
CREATE TABLE IF NOT EXISTS TagIndex
(
    `tag`    VARCHAR(128) NOT NULL,
    `caseId` INT UNSIGNED NOT NULL,
    `weight` FLOAT        NOT NULL,
    FOREIGN KEY (`caseId`) REFERENCES Cases (`id`)
) CHAR SET utf8;
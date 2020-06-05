CREATE TABLE IF NOT EXISTS Cases # 案件数据库
(
    `id`     INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `judges` TEXT         NULL,     # 法官 考虑到一对多关系 倒排索引需单独建表
    `laws`   TEXT         NULL,     # 法律 考虑到一对多关系 倒排索引需单独建表
    `tags`   TEXT         NULL,     # 标签 考虑到一对多关系 倒排索引需单独建表
    `detail` LONGTEXT     NOT NULL, # 案情 xml 中的全文部分
    `tree`   LONGTEXT     NOT NULL, # xml 的树形结构
    PRIMARY KEY (`id` ASC)
) CHAR SET utf8;
CREATE TABLE IF NOT EXISTS WordIndex # 词倒排索引数据库, word -> [](case, weight)
(
    `word`   VARCHAR(64)  NOT NULL,
    `caseId` INT UNSIGNED NOT NULL, # Foreign key 指向案件
    `weight` FLOAT        NOT NULL, # 将利用出现词的个数等等进行简单的计算
    PRIMARY KEY (`word`, `caseId`),
    INDEX `word_weight_idx` (`weight`)
) CHAR SET utf8;
CREATE TABLE IF NOT EXISTS JudgeIndex
(
    `judge`  VARCHAR(64)  NOT NULL,
    `caseId` INT UNSIGNED NOT NULL, # Foreign key 指向案件
    `weight` FLOAT        NOT NULL, # 将利用出现法官出现顺序等等进行简单的计算
    PRIMARY KEY (`judge`, `caseId`),
    INDEX `judge_weight_idx` (`weight`)
) CHAR SET utf8;
CREATE TABLE IF NOT EXISTS LawIndex
(
    `law`    VARCHAR(512) NOT NULL,
    `caseId` INT UNSIGNED NOT NULL, # Foreign key 指向案件
    `weight` FLOAT        NOT NULL, # 将利用出现法律的个数等等进行简单的计算
    PRIMARY KEY (`law`, `caseId`),
    INDEX `law_weight_idx` (`weight`)
) CHAR SET utf8;
CREATE TABLE IF NOT EXISTS TagIndex
(
    `tag`    VARCHAR(128) NOT NULL,
    `caseId` INT UNSIGNED NOT NULL, # Foreign key 指向案件
    `weight` FLOAT        NOT NULL, # 一些重要的标签和 jieba 分词的 topK 功能
    PRIMARY KEY (`tag`, `caseId`),
    INDEX `tag_weight_idx` (`weight`)
) CHAR SET utf8;

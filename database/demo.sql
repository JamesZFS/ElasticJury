# # 让mysql查一次，就查出这么多个词的caseId，并且把相同caseId的weight求和:
# select caseId, sum(weight) as weight
# from WordIndex
# where word in ('民事', '调解书', '内容')
# group by caseId
# order by weight desc
# limit 10;

CREATE TEMPORARY TABLE IF NOT EXISTS WordWeight
(
    word   VARCHAR(64) NOT NULL, # 用户输入的查询词
    weight FLOAT       NOT NULL, # 用户输入的词的权重（idf或者输入词的次数）
    PRIMARY KEY (word)           # 一对一映射
) CHAR SET utf8;

truncate WordWeight;
insert into WordWeight (word, weight)
VALUES ('民事', 1.0),
       ('调解书', 2.0),
       ('内容', 3.0);

# 按照WordWeight表中的词权重进行检索、求和、排序
select a.caseId as caseId, sum(a.weight * b.weight) as weight
from WordIndex a,
     WordWeight b
where a.word = b.word
  and a.word in ('民事', '调解书', '内容')
group by caseId
order by weight desc
limit 10;

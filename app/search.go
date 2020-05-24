package app

import (
	. "ElasticJury/app/common"
	"ElasticJury/app/natural"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Handle search case id by words/tags/laws/judges
// Method: POST
//
// Queries:
// 		word: "调解,当事人,..." separated by ','
// 		tags: "民事案件,一审案件,..." separated by ','
// 		laws: "中华人民共和国民法通则,《中华人民共和国担保法》,..." separated by ','  quoted by '《》' or not
// 		judges: "黄琴英,高原,..." separated by ','
//
// Params(json):
//      misc: miscellaneous searching field, a text representing a case description. This field will be automatically
//			divided into the four fields above for searching.
//
func (db database) makeSearchHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		// Parse queries and params:
		words := natural.PreprocessWords(strings.Split(context.Query("word"), ","))
		tags := FilterStrs(strings.Split(context.Query("tag"), ","), NotWhiteSpace)
		laws := FilterStrs(strings.Split(context.Query("law"), ","), NotWhiteSpace)
		judges := FilterStrs(strings.Split(context.Query("judge"), ","), NotWhiteSpace)
		var json struct {
			Misc string `json:"misc" form:"misc"`
		}
		if err := context.BindJSON(&json); err != nil && err != io.EOF { // parsing from post data
			_ = context.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if NotWhiteSpace(json.Misc) {
			// Parse misc text and output it into the four fields
			words_, tags_, laws_, judges_ := natural.ParseFullText(json.Misc)
			words = append(words, words_...)
			tags = append(tags, tags_...)
			laws = append(laws, laws_...)
			judges = append(judges, judges_...)
		}

		// Perform searching for each field:
		var (
			result    searchResultSet
			newResult searchResultSet
			err       error
		)
		if len(words) > 0 {
			wordWeightMap := map[string]float32{}
			for _, word := range words {
				wordWeightMap[word] = 1.0 // TODO, maybe use
			}
			result, err = db.searchCaseIdsByWordWeightMap(wordWeightMap, WordSearchLimit)
			if err != nil {
				panic(err)
			}
			if len(result) == 0 { // early stop with empty response
				context.JSON(http.StatusOK, emptyResponse)
				return
			}
		}
		if len(tags) > 0 { //  todo: maybe we do refined search in word search result set
			newResult, err = db.searchCaseIdsByTag(tags)
			if err != nil {
				panic(err)
			}
			result = newResult.merge(result)
			if len(result) == 0 { // early stop with empty response
				context.JSON(http.StatusOK, emptyResponse)
				return
			}
		}
		if len(laws) > 0 {
			newResult, err = db.searchCaseIdsByLaw(laws)
			if err != nil {
				panic(err)
			}
			result = newResult.merge(result)
			if len(result) == 0 { // early stop with empty response
				context.JSON(http.StatusOK, emptyResponse)
				return
			}
		}
		if len(judges) > 0 {
			newResult, err = db.searchCaseIdsByJudge(judges)
			if err != nil {
				panic(err)
			}
			result = newResult.merge(result)
		}
		if result == nil {
			result = searchResultSet{}
		}

		context.JSON(http.StatusOK, result.sortMapByValue().toResponse())
	}
}

func (db database) makeCaseInfoHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		idQuery := context.Query("id")
		ids := strings.Split(idQuery, ",")
		for _, id := range ids { // id check (avoid injection e.t.c.)
			if _, err := strconv.ParseInt(id, 10, 64); err != nil {
				context.String(http.StatusBadRequest, "Bad `id` query \"%id\".", id)
				return
			}
		}
		rows, err := db.Query(fmt.Sprintf("SELECT id, judges, laws, tags, detail FROM Cases WHERE id IN (%s) ORDER BY FIELD(id, %s)", idQuery, idQuery))
		if err != nil {
			panic(err)
		}
		result := make([]gin.H, 0, len(ids))
		for rows.Next() {
			var (
				id                         int
				judges, laws, tags, detail string
			)
			if err := rows.Scan(&id, &judges, &laws, &tags, &detail); err != nil {
				panic(err)
			}
			result = append(result, gin.H{
				"id":     id,
				"judges": judges,
				"laws":   laws,
				"tags":   tags,
				"detail": detail,
			})
		}
		context.JSON(http.StatusOK, result)
	}
}

// Or mode
func (db database) searchBy(querySql string, keys []string) (searchResultSet, error) {
	result := searchResultSet{}
	for _, key := range keys { // query each caseId in `WordIndex` table
		rows, err := db.Query(querySql, key)
		if err != nil {
			return nil, err
		}
		for rows.Next() { // append new case to the result set
			var (
				caseId int
				weight float32
			)
			if err := rows.Scan(&caseId, &weight); err != nil {
				return nil, err
			}
			result[caseId] += weight
		}
	}
	return result, nil
}

// Old method
func (db database) searchCaseIdsByWord(words []string) (searchResultSet, error) {
	// language=MySQL
	return db.searchBy("SELECT `caseId`, `weight` FROM WordIndex WHERE `word` = ?", words)
}

// Input: word(user input) -> weight(input word count or idf)
func (db database) searchCaseIdsByWordWeightMap(wordWeightMap map[string]float32, limit int) (result searchResultSet, err error) {
	// Create table
	tableId := time.Now().UnixNano()
	createTable := fmt.Sprintf(`
		CREATE TABLE WordWeight%d
		(
			word   VARCHAR(64) NOT NULL, # 用户输入的查询词
			weight FLOAT       NOT NULL, # 用户输入的词的权重（idf或者输入词的次数）
			PRIMARY KEY (word)           # 一对一映射
		) CHAR SET utf8;`, tableId)
	if _, err = db.Exec(createTable); err != nil {
		return nil, err
	}

	// Drop after finish
	defer func() {
		drop := fmt.Sprintf(`DROP TABLE WordWeight%d`, tableId)
		if _, err = db.Exec(drop); err != nil { // Overwrite return value
			result = nil
		}
	}()

	// Insert items
	items := make([]string, 0, len(wordWeightMap))
	for word, weight := range wordWeightMap {
		items = append(items, fmt.Sprintf("('%s',%f)", word, weight))
	}
	insert := fmt.Sprintf(`INSERT INTO WordWeight%d (word, weight) VALUES %s`, tableId, strings.Join(items, ","))
	if _, err = db.Exec(insert); err != nil {
		return nil, err
	}

	// Search
	query := fmt.Sprintf(`
		select a.caseId as caseId, sum(a.weight * b.weight) as weight
		from WordIndex a, WordWeight%d b
		where a.word = b.word
		group by caseId order by weight desc limit %d`, tableId, limit)
	var rows *sql.Rows
	rows, err = db.Query(query)
	if err != nil {
		return nil, err
	}
	result = searchResultSet{}
	for rows.Next() {
		var (
			caseId int
			weight float32
		)
		if err = rows.Scan(&caseId, &weight); err != nil {
			return nil, err
		}
		result[caseId] = weight
	}

	return result, err
}

// todo: maybe we limit the tag/law/judge result count? like we do in words
func (db database) searchCaseIdsByTag(tags []string) (searchResultSet, error) {
	// language=MySQL
	return db.searchBy("SELECT `caseId`, `weight` FROM TagIndex WHERE `tag` = ?", tags)
}

func (db database) searchCaseIdsByLaw(laws []string) (searchResultSet, error) {
	// language=MySQL
	return db.searchBy("SELECT `caseId`, `weight` FROM LawIndex WHERE `law` = ?", laws)
}

func (db database) searchCaseIdsByJudge(judges []string) (searchResultSet, error) {
	// language=MySQL
	return db.searchBy("SELECT `caseId`, `weight` FROM JudgeIndex WHERE `judge` = ?", judges)
}

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
		// Parse params:
		words := natural.PreprocessWords(strings.Split(context.Query("word"), ","))
		filters := []Filter {
			BuildFilter("TagIndex", "tag", context.Query("tag")),
			BuildFilter("LawIndex", "law", context.Query("law")),
			BuildFilter("JudgeIndex", "judge", context.Query("judge")),
		}

		// Parse content
		var json struct {
			Misc string `json:"misc" form:"misc"`
		}
		if err := context.BindJSON(&json); err != nil && err != io.EOF { // parsing from post data
			_ = context.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if NotWhiteSpace(json.Misc) {
			// Parse misc text and output it into the four fields
			words = append(words, natural.ParseFullText(json.Misc)...)
		}

		// Perform searching
		wordWeightMap := map[string]float32{}
		for _, word := range words {
			wordWeightMap[word] = 1.0 // TODO: maybe use
		}
		result, err := db.searchCaseIdsByWordWeightMap(wordWeightMap, filters, WordSearchLimit)
		if err != nil {
			panic(err)
		}

		// Return
		context.JSON(http.StatusOK, result.sortMapByValue().toResponse())
	}
}

func (db database) makeCaseInfoHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		idQuery := context.Query("id")
		ids := strings.Split(idQuery, ",")

		// Checks
		for _, id := range ids {
			if _, err := strconv.ParseInt(id, 10, 64); err != nil {
				context.String(http.StatusBadRequest, "Bad `id` query \"%id\".", id)
				return
			}
		}

		// Query
		rows, err := db.Query(fmt.Sprintf(`
			SELECT id, judges, laws, tags, detail
			FROM Cases WHERE id IN (%s)
			ORDER BY FIELD(id, %s)`, idQuery, idQuery))
		if err != nil {
			panic(err)
		}

		// Return
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

// Input: word(user input) -> weight(input word count or idf)
func (db database) searchCaseIdsByWordWeightMap(wordWeightMap map[string]float32, filters []Filter, limit int) (result searchResultSet, err error) {
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
		return searchResultSet{}, err
	}

	// Drop after finish
	defer func() {
		drop := fmt.Sprintf(`DROP TABLE WordWeight%d`, tableId)
		if _, err = db.Exec(drop); err != nil { // Overwrite return value
			result = searchResultSet{}
		}
	}()

	// Insert items
	items := make([]string, 0, len(wordWeightMap))
	for word, weight := range wordWeightMap {
		items = append(items, fmt.Sprintf("('%s',%f)", word, weight))
	}
	insert := fmt.Sprintf(`INSERT INTO WordWeight%d (word, weight) VALUES %s`, tableId, strings.Join(items, ","))
	if _, err = db.Exec(insert); err != nil {
		return searchResultSet{}, err
	}

	// Convert filters into SQL conditions
	tables, conditions, entryIndex := "", "", 'c'
	for _, filter := range filters {
		if len(filter.Conditions) > 0 {
			tables += fmt.Sprintf(", %s %c", filter.TableName, entryIndex)
			orExpr := GetOrExpr(entryIndex, filter.FieldName, filter.Conditions)
			conditions += fmt.Sprintf(" AND a.caseId=%c.caseId AND (%s)", entryIndex, orExpr)
			entryIndex ++
		}
	}

	// Search
	query := fmt.Sprintf(`
		SELECT a.caseId AS caseId, sum(a.weight * b.weight) AS weight
		FROM WordIndex a, WordWeight%d b%s
		WHERE a.word = b.word%s
		GROUP BY caseId ORDER BY weight DESC LIMIT %d`, tableId, tables, conditions, limit)
	println(query)
	var rows *sql.Rows
	rows, err = db.Query(query)
	if err != nil {
		return searchResultSet{}, err
	}
	result = searchResultSet{}
	for rows.Next() {
		var (
			caseId int
			weight float32
		)
		if err = rows.Scan(&caseId, &weight); err != nil {
			return searchResultSet{}, err
		}
		result[caseId] = weight
	}

	return result, err
}
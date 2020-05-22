package app

import (
	. "ElasticJury/app/common"
	"ElasticJury/app/natural"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type (
	searchResultSet map[int]float32 // caseId -> weight
)

const (
	modeAnd = false // Intersect search results, default
	modeOr  = true  // Union search results
)

var (
	emptySearchResponse = searchResultSet{}.toResponse()
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

		var mode bool // fixme: this is deprecated now
		if context.Query("mode") == "OR" {
			mode = modeOr
		} else {
			mode = modeAnd
		}

		// Perform searching for each field:
		var (
			result    searchResultSet
			newResult searchResultSet
			err       error
		)
		if len(words) > 0 {
			result, err = db.searchCaseIdsByWord(words, mode)
			if err != nil {
				panic(err)
			}
			if len(result) == 0 { // early stop with empty response
				context.JSON(http.StatusOK, emptySearchResponse)
				return
			}
		}
		if len(tags) > 0 {
			newResult, err = db.searchCaseIdsByTag(tags, mode)
			if err != nil {
				panic(err)
			}
			result = newResult.merge(result)
			if len(result) == 0 { // early stop with empty response
				context.JSON(http.StatusOK, emptySearchResponse)
				return
			}
		}
		if len(laws) > 0 {
			newResult, err = db.searchCaseIdsByLaw(laws, mode)
			if err != nil {
				panic(err)
			}
			result = newResult.merge(result)
			if len(result) == 0 { // early stop with empty response
				context.JSON(http.StatusOK, emptySearchResponse)
				return
			}
		}
		if len(judges) > 0 {
			newResult, err = db.searchCaseIdsByJudge(judges, mode)
			if err != nil {
				panic(err)
			}
			result = newResult.merge(result)
		}
		if result == nil {
			result = searchResultSet{}
		}

		context.JSON(http.StatusOK, result.toResponse())
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

func (db database) searchBy(querySql string, keys []string, mode bool) (searchResultSet, error) {
	result := searchResultSet{}
	for i, key := range keys { // query each key in `WordIndex` table
		newResult := searchResultSet{}
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
			if mode == modeAnd {
				if i == 0 {
					newResult[caseId] = weight
				} else {
					oldWeight, contains := result[caseId]
					if contains {
						newResult[caseId] = oldWeight + weight
					}
				}
			} else { // modeOr
				result[caseId] += weight
			}
		}
		if mode == modeAnd {
			result = newResult
			if len(result) == 0 { // early stop if empty
				return result, nil
			}
		}
	}
	return result, nil
}

func (db database) searchCaseIdsByWord(words []string, mode bool) (searchResultSet, error) {
	// language=MySQL
	return db.searchBy("SELECT `caseId`, `weight` FROM WordIndex WHERE `word` = ?", words, mode)
}

func (db database) searchCaseIdsByTag(tags []string, mode bool) (searchResultSet, error) {
	// language=MySQL
	return db.searchBy("SELECT `caseId`, `weight` FROM TagIndex WHERE `tag` = ?", tags, mode)
}

func (db database) searchCaseIdsByLaw(laws []string, mode bool) (searchResultSet, error) {
	// language=MySQL
	return db.searchBy("SELECT `caseId`, `weight` FROM LawIndex WHERE `law` = ?", laws, mode)
}

func (db database) searchCaseIdsByJudge(judges []string, mode bool) (searchResultSet, error) {
	// language=MySQL
	return db.searchBy("SELECT `caseId`, `weight` FROM JudgeIndex WHERE `judge` = ?", judges, mode)
}

// Intersect search results, nil set stands for **full set**
func (s searchResultSet) merge(t searchResultSet) searchResultSet {
	if s == nil {
		return t
	}
	if t == nil {
		return s
	}
	if len(t) < len(s) {
		return t.merge(s)
	}
	// Assume len(s) <= len(t)
	res := searchResultSet{}
	for id, w1 := range s {
		if w2, contains := t[id]; contains {
			res[id] = w1 + w2 // TODO maybe other operations
		}
	}
	return res
}

// To http response body
func (s searchResultSet) toResponse() gin.H {
	return gin.H{
		"count":  len(s),
		"result": s,
	}
}

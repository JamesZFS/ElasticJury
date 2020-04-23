package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

func (db database) makeSearchHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		s := searchResultSet{}
		s[1] = 1.1
		// Parse queries:
		words := PreprocessWords(strings.Split(context.Query("word"), ","))
		tags := FilterStrs(strings.Split(context.Query("tag"), ","), NotWhiteSpace)
		laws := FilterStrs(strings.Split(context.Query("law"), ","), NotWhiteSpace)
		judges := FilterStrs(strings.Split(context.Query("judge"), ","), NotWhiteSpace)
		var mode bool
		if context.Query("mode") == "OR" {
			mode = modeOr
		} else {
			mode = modeAnd
		}

		// Performance searching for each field:
		var (
			result    searchResultSet
			newResult searchResultSet
			err       error
		)
		if len(words) > 0 {
			result, err = db.searchCaseIdsByWord(words, mode)
			if err != nil {
				context.Status(http.StatusInternalServerError)
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
				context.Status(http.StatusInternalServerError)
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
				context.Status(http.StatusInternalServerError)
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
				context.Status(http.StatusInternalServerError)
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

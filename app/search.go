package app

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

//type searchResultItem struct {
//	caseId int
//	weight float32
//}

type (
	searchResultSet map[int]float32 // caseId -> weight
	caseIdSearcher  func(db *sql.DB, keys []string, mode bool) (searchResultSet, error)
)

const (
	modeAnd = false // Intersect search results, default
	modeOr  = true  // Union search results
)

// Since go has no procedural macro, we can only do this in runtime...
var (
	searchCaseIdsByWord  = makeSearcher("WordIndex", "word")
	searchCaseIdsByTag   = makeSearcher("TagIndex", "tag")
	searchCaseIdsByLaw   = makeSearcher("LawIndex", "law")
	searchCaseIdsByJudge = makeSearcher("JudgeIndex", "judge")
	emptySearchResponse  = gin.H{
		"count":  0,
		"result": searchResultSet{},
	}
)

func MakeSearchHandler(db *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
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
		var (
			result    searchResultSet
			newResult searchResultSet
			err       error
		)
		if len(words) > 0 {
			result, err = searchCaseIdsByWord(db, words, mode)
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
			newResult, err = searchCaseIdsByTag(db, tags, mode)
			if err != nil {
				context.Status(http.StatusInternalServerError)
				panic(err)
			}
			result = mergeSearchResult(newResult, result)
			if len(result) == 0 { // early stop with empty response
				context.JSON(http.StatusOK, emptySearchResponse)
				return
			}
		}
		if len(laws) > 0 {
			newResult, err = searchCaseIdsByLaw(db, laws, mode)
			if err != nil {
				context.Status(http.StatusInternalServerError)
				panic(err)
			}
			result = mergeSearchResult(newResult, result)
			if len(result) == 0 { // early stop with empty response
				context.JSON(http.StatusOK, emptySearchResponse)
				return
			}
		}
		if len(judges) > 0 {
			newResult, err = searchCaseIdsByJudge(db, judges, mode)
			if err != nil {
				context.Status(http.StatusInternalServerError)
				panic(err)
			}
			result = mergeSearchResult(newResult, result)
		}

		context.JSON(http.StatusOK, gin.H{
			"count":  len(result),
			"result": result,
		})
	}
}

// Return a case id search that search in `tableName`, using the WHERE clause `keyName` = ?
func makeSearcher(tableName string, keyName string) caseIdSearcher {
	querySql := fmt.Sprintf("SELECT `caseId`, `weight` FROM %s WHERE `%s` = ?", tableName, keyName)
	return func(db *sql.DB, keys []string, mode bool) (searchResultSet, error) {
		result := searchResultSet{}
		for i, key := range keys { // query each key in `WordIndex` table
			//var newResult searchResultSet
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
}

// Intersect search results, nil set stands for **full set**
func mergeSearchResult(set1, set2 searchResultSet) searchResultSet {
	if set1 == nil {
		return set2
	}
	if set2 == nil {
		return set1
	}
	if len(set2) < len(set1) {
		return mergeSearchResult(set2, set1)
	}
	// Assume len(set1) <= len(set2)
	res := searchResultSet{}
	for id, w1 := range set1 {
		if w2, contains := set2[id]; contains {
			res[id] = w1 + w2 // TODO maybe other operations
		}
	}
	return res
}

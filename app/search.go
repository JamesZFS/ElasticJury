package app

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

//type searchResultItem struct {
//	caseId int
//	weight float32
//}

type searchResultSet map[int]float32 // caseId -> weight

const (
	modeAnd = false // Intersect search results, default
	modeOr  = true  // Union search results
)

func MakeSearchHandler(db *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		words := strings.Split(context.Query("words"), ",")
		var mode bool
		if context.Query("mode") == "OR" {
			mode = modeOr
		} else {
			mode = modeAnd
		}
		result, err := searchCaseIdsByWords(db, words, mode)
		if err != nil {
			context.Status(http.StatusInternalServerError)
			panic(err)
		}

		context.JSON(http.StatusOK, gin.H{
			"count":  len(result),
			"result": result,
		})
	}
}

func searchCaseIdsByWords(db *sql.DB, words []string, mode bool) (searchResultSet, error) {
	result := searchResultSet{}
	for i, word := range words { // query each word in `WordIndex` table
		//var newResult searchResultSet
		newResult := searchResultSet{}
		word = preprocessWord(word)
		rows, err := db.Query("SELECT `caseId`, `weight` FROM WordIndex WHERE `word` = ?", word)
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
		}
	}
	return result, nil
}

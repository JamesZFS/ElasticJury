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
		// Parse content
		var json struct {
			Misc string `json:"misc" form:"misc"`
			Tag string `json:"tag" form:"tag"`
			Law string `json:"law" form:"law"`
			Judge string `json:"judge" form:"judge"`
		}
		if err := context.BindJSON(&json); err != nil && err != io.EOF { // parsing from post data
			_ = context.AbortWithError(http.StatusBadRequest, err)
			return
		}

		fmt.Printf("[Search] Got request:\n")
		fmt.Printf("[Search]  > Misc: %s\n", json.Misc)
		fmt.Printf("[Search]  > Tag: %s\n", json.Tag)
		fmt.Printf("[Search]  > Law: %s\n", json.Law)
		fmt.Printf("[Search]  > Judge: %s\n", json.Judge)

		var words Conditions
		if NotWhiteSpace(json.Misc) {
			words = natural.ParseFullText(json.Misc)
		}
		tags := natural.PreprocessWords(strings.Split(json.Tag, ","))
		laws := natural.PreprocessWords(strings.Split(json.Law, ","))
		judges := natural.PreprocessWords(strings.Split(json.Judge, ","))

		// Params
		params := []Param{
			BuildParam("WordIndex", "word", words),
			BuildParam("TagIndex", "tag", MakeDefaultConditions(tags)),
			BuildParam("LawIndex", "law", MakeDefaultConditions(laws)),
			BuildParam("JudgeIndex", "judge", MakeDefaultConditions(judges)),
		}

		// Perform searching
		result, err := db.searchCaseIds(params)
		if err != nil {
			if err, castSuccess := err.(EmptyParamErr); castSuccess {
				context.Data(http.StatusOK, "binary", []byte{})
				fmt.Printf("[Request] Failed with error: %s\n", err.Error())
				return
			}
			panic(err) // unknown
		}

		// Return
		fmt.Printf("[Search] Reply with %d cases\n", len(result))
		context.Data(http.StatusOK, "binary", result.ToByteArray())
	}
}

func (db database) makeCaseInfoHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		// Parse content
		var json struct {
			Id string `json:"id" form:"id"`
		}
		if err := context.BindJSON(&json); err != nil && err != io.EOF { // parsing from post data
			_ = context.AbortWithError(http.StatusBadRequest, err)
			return
		}
		idQuery := json.Id
		ids := strings.Split(idQuery, ",")
		fmt.Printf("[CaseInfo] Got request:\n")
		fmt.Printf("[CaseInfo]  > IDs: %s\n", idQuery)

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
		fmt.Printf("[CaseInfo] Reply with %d infos\n", len(result))
		context.JSON(http.StatusOK, result)
	}
}

func (db database) makeCaseDetailHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		id, err := strconv.Atoi(context.Param("id"))
		fmt.Printf("[Detail] Got request:\n")
		fmt.Printf("[Detail]  > ID: %d\n", id)
		if id <= 0 || err != nil {
			_ = context.AbortWithError(http.StatusBadRequest, err)
			return
		}
		var (
			judges, laws, tags, detail, tree string
		)
		if err := db.QueryRow("SELECT judges, laws, tags, detail, tree FROM Cases WHERE id = ?", id).
			Scan(&judges, &laws, &tags, &detail, &tree); err != nil {
			panic(err)
		}
		fmt.Printf("[Detail] Reply successfully\n")
		context.JSON(http.StatusOK, gin.H{
			"id":     id,
			"judges": judges,
			"laws":   laws,
			"tags":   tags,
			"detail": detail,
			"tree":   tree,
		})
	}
}

func (db database) searchCaseIds(params []Param) (result ResultList, err error) {
	// Create table
	tableId := time.Now().UnixNano()
	createTable := fmt.Sprintf(`
		CREATE TABLE Weights%d
		(
			item   VARCHAR(512) NOT NULL,  # 首要的检索条件 
			weight FLOAT        NOT NULL,  # 用户输入的词的权重（idf或者输入词的次数）
			PRIMARY KEY (item)             # 一对一映射
		) CHAR SET utf8;`, tableId)
	if _, err = db.Exec(createTable); err != nil {
		return ResultList{}, err
	}

	// Drop after finish
	defer func() {
		drop := fmt.Sprintf(`DROP TABLE Weights%d`, tableId)
		if _, errDrop := db.Exec(drop); err == nil && errDrop != nil {
			result, err = ResultList{}, errDrop
		}
	}()

	// Convert params into SQL conditions
	tables, conditions, entryIndex, first := "", "", 'b', true
	for _, param := range params {
		if len(param.Conditions) > 0 {
			tables += fmt.Sprintf(", %s %c", param.TableName, entryIndex)
			if first {
				first = false
				conditions += fmt.Sprintf("a.item = %c.%s", entryIndex, param.FieldName)

				// Insert items
				var items []string
				for _, condition := range param.Conditions {
					items = append(items, fmt.Sprintf("('%s',%f)", condition.Item, condition.Weight))
				}
				insert := fmt.Sprintf(`INSERT INTO Weights%d (item, weight) VALUES %s`, tableId, strings.Join(items, ","))
				if _, err = db.Exec(insert); err != nil {
					return ResultList{}, err
				}
			} else {
				orExpr := GetOrExpr(entryIndex, param.FieldName, param.Conditions)
				conditions += fmt.Sprintf(" AND b.caseId=%c.caseId AND (%s)", entryIndex, orExpr)
			}
			entryIndex++
		}
	}
	if first {
		return nil, EmptyParamErr{}
	}

	// Search
	var limit string
	if // noinspection ALL
	SearchLimit > 0 {
		limit = fmt.Sprintf("LIMIT %d", SearchLimit)
	} else {
		limit = ""
	}
	query := fmt.Sprintf(`
		SELECT b.caseId AS caseId, sum(b.weight * a.weight) AS weight
		FROM Weights%d a%s
		WHERE %s
		GROUP BY caseId ORDER BY weight DESC %s`, tableId, tables, conditions, limit)
	var rows *sql.Rows
	rows, err = db.Query(query)
	if err != nil {
		return ResultList{}, err
	}
	result = ResultList{}
	for rows.Next() {
		var (
			caseId int
			weight float32
		)
		if err = rows.Scan(&caseId, &weight); err != nil {
			return ResultList{}, err
		}
		if weight > SearchFilter {
			result = append(result, int32(caseId))
		}
	}

	return result, err
}

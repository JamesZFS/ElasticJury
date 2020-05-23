package app

import (
	. "ElasticJury/app/common"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math"
	"net/http/httptest"
	"strings"
	"testing"
)

const epsilon = 1e-3

func TestApp(t *testing.T) {
	app := appPrologue()
	defer appEpilogue(app)
	var searchTests = []struct {
		name string
		path string // request
		want gin.H  // response
	}{
		{
			name: "Search for nil",
			path: "/search",
			want: gin.H{
				"count":  0,
				"result": map[string]interface{}{},
			},
		},
		{
			name: "Search for word 1",
			path: "/search?word=世界,你好",
			want: gin.H{
				"count": 1,
				"result": map[string]interface{}{
					"1": 1.3,
				},
			},
		},
		{
			name: "Search for word 2",
			path: "/search?word=世界,你好,&mode=OR",
			want: gin.H{
				"count": 3,
				"result": map[string]interface{}{
					"1": 1.3,
					"2": 0.3,
					"3": 0.3,
				},
			},
		},
		{
			name: "Search for word not found",
			path: "/search?word=世界,啥玩意,&mode=AND",
			want: gin.H{
				"count":  0,
				"result": map[string]interface{}{},
			},
		},
		{
			name: "Search for tag",
			path: "/search?tag=猥亵",
			want: gin.H{
				"count": 2,
				"result": map[string]interface{}{
					"2": 0.6,
					"3": 0.6,
				},
			},
		},
		{
			name: "Search for law",
			path: "/search?law=宪法",
			want: gin.H{
				"count": 2,
				"result": map[string]interface{}{
					"1": 0.9,
					"2": 0.9,
				},
			},
		},
		{
			name: "Search for judge",
			path: "/search?judge=王五",
			want: gin.H{
				"count": 1,
				"result": map[string]interface{}{
					"2": 1.0,
				},
			},
		},
		{
			name: "Search misc 1",
			path: "/search?judge=王五&law=宪法",
			want: gin.H{
				"count": 1,
				"result": map[string]interface{}{
					"2": 1.9,
				},
			},
		},
		{
			name: "Search misc 2",
			path: "/search?tag=诈骗&law=诉讼法",
			want: gin.H{
				"count":  0,
				"result": map[string]interface{}{},
			},
		},
		{
			name: "Search misc 3",
			path: "/search?tag=猥亵&law=宪法,诉讼法&word=某人&mode=OR",
			want: gin.H{
				"count": 2,
				"result": map[string]interface{}{
					"2": 0.6 + 0.9 + 0.3,
					"3": 0.6 + 1.0 + 0.3,
				},
			},
		},
	}
	t.Run("Ping test", func(t *testing.T) {
		got := string(Request("GET", "/ping", app.Engine))
		assert.Equal(t, "pong", strings.ToLower(got))
	})

	for _, tt := range searchTests {
		t.Run(tt.name, func(t *testing.T) {
			body := Request("POST", tt.path, app.Engine)
			var got gin.H
			assert.Nilf(t, json.Unmarshal(body, &got), "Fail to parse response body %s", body)
			assert.Equal(t, tt.want["count"], int(got["count"].(float64)))
			assert.Truef(t, resultApproxEqual(got["result"].(map[string]interface{}), tt.want["result"].(map[string]interface{})),
				"Search result compare failed!\ngot  %v\nwant %v", got, tt.want)
		})
	}
}

// Request 根据特定请求uri，发起get/post请求返回响应
func Request(method string, path string, router *gin.Engine) []byte {
	// 构造get请求
	req := httptest.NewRequest(method, path, nil)
	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	router.ServeHTTP(w, req)

	// 提取响应
	result := w.Result()

	// 读取响应body
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}

	if err := result.Body.Close(); err != nil {
		panic(err)
	}
	return body
}

func appPrologue() *App {
	dbPrologue()
	password := GetEnvVar("PASSWORD", "")
	app := NewApp(TestDatabaseName, password)
	return app
}

func appEpilogue(app *App) {
	dbEpilogue(app.db)
}

func resultApproxEqual(s1, s2 map[string]interface{}) bool {
	for i, f1 := range s1 {
		f2, ok := s2[i]
		if !ok {
			return false
		}
		if math.Abs(f1.(float64)-f2.(float64)) > epsilon {
			return false
		}
	}
	for i, f1 := range s2 {
		f2, ok := s1[i]
		if !ok {
			return false
		}
		if math.Abs(f1.(float64)-f2.(float64)) > epsilon {
			return false
		}
	}
	return true
}

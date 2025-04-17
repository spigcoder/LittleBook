package test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spigcoder/LittleBook/webook/internal/repository/dao"
	"github.com/spigcoder/LittleBook/webook/internal/web/ijwt"
	"github.com/spigcoder/LittleBook/webook/ioc"
	"github.com/spigcoder/LittleBook/webook/startup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (s *ArticleTestSuite) SetupSuite() {
	s.server = gin.Default()
	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("claims",
			&ijwt.UserClaims{
				Uid: 123,
			})
	})
	ioc.InitViper()
	artHdl := startup.InitArticleHandler()
	s.db = startup.InitTestDB()
	artHdl.RegisterRoutes(s.server)
}

func (s *ArticleTestSuite) TearDownSuite() {
	s.db.Exec("truncate table articles")
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCase := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		//预期输入
		article Article

		//预期输出
		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "编辑成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art dao.Article
				err := s.db.Where("id = ?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.CTime > 0)
				assert.True(t, art.UTime > 0)
				art.CTime = 0
				art.UTime = 0
				assert.Equal(t, dao.Article{
					Id:       1,
					Title:    "标题",
					Content:  "内容",
					AuthorId: 123,
				}, art)
			},
			article: Article{
				Title:   "标题",
				Content: "内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: http.StatusOK,
				Msg:  "发布成功",
				Data: 1,
			},
		},
		{
			name: "编辑更改",
			before: func(t *testing.T) {
				//首先插入数据
				err := s.db.Create(&dao.Article{
					Id:       1,
					Title:    "标题",
					Content:  "内容",
					AuthorId: 123,
					CTime:    123,
					UTime:    123,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art dao.Article
				err := s.db.Where("id = ?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.CTime > 0)
				assert.True(t, art.UTime > 123)
				art.UTime = 0
				assert.Equal(t, dao.Article{
					Id:       1,
					Title:    "新的标题",
					Content:  "新的内容",
					AuthorId: 123,
					CTime:    123,
				}, art)
			},
			article: Article{
				Title:   "新的标题",
				Content: "新的内容",
				Id:      1,
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: http.StatusOK,
				Msg:  "发布成功",
				Data: 1,
			},
		},
		{
			name: "有人错误的更改你的信息",
			before: func(t *testing.T) {
				//首先插入数据
				err := s.db.Create(&dao.Article{
					Id:       4,
					Title:    "标题",
					Content:  "内容",
					AuthorId: 234,
					CTime:    123,
					UTime:    123,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art dao.Article
				err := s.db.Where("id = ?", 4).First(&art).Error
				assert.NoError(t, err)
				art.UTime = 0
				assert.Equal(t, dao.Article{
					Id:       4,
					Title:    "标题",
					Content:  "内容",
					AuthorId: 234,
					UTime:    123,
					CTime:    123,
				}, art)
			},
			article: Article{
				Title:   "标题",
				Content: "内容",
				Id:      4,
			},
			wantCode: http.StatusInternalServerError,
			wantRes: Result[int64]{
				Code: http.StatusInternalServerError,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			//制造数据
			//tc.before(t)
			//defer tc.after(t)
			reqBody, err := json.Marshal(tc.article)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewReader(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			//执行请求
			s.server.ServeHTTP(resp, req)
			//交验结果
			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != http.StatusOK {
				return
			}
			var res Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

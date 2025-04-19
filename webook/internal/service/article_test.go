package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/spigcoder/LittleBook/webook/internal/domain"
	"github.com/spigcoder/LittleBook/webook/internal/repository"
	"github.com/spigcoder/LittleBook/webook/internal/repository/article"
	artRepoMocks "github.com/spigcoder/LittleBook/webook/internal/repository/article/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArticleService_Publish(t *testing.T) {
	testCase := []struct {
		name    string
		mock    func(ctl *gomock.Controller) (article.ArticleAuthorRepository, repository.ArticleRepository)
		Article domain.Article
		wantID  int64
		wantErr error
	}{
		{
			name: "发布文章成功",
			Article: domain.Article{
				Title:   "标题",
				Content: "内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantID:  1,
			wantErr: nil,
			mock: func(ctl *gomock.Controller) (article.ArticleAuthorRepository, repository.ArticleRepository) {
				author := artRepoMocks.NewMockArticleAuthorRepository(ctl)
				author.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				reader := artRepoMocks.NewMockArticleReaderRepository(ctl)
				reader.EXPECT().Create(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123},
				}).Return(int64(1), nil)
				return author, reader
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			svc := NewArticleService(tc.mock(ctl))
			id, err := svc.Publish(context.Background(), tc.Article)
			assert.Equal(t, id, tc.wantID)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}

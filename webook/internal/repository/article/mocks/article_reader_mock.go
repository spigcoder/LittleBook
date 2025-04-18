// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/article/article_reader.go

// Package artRepoMocks is a generated GoMock package.
package artRepoMocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/spigcoder/LittleBook/webook/internal/domain"
)

// MockArticleReaderRepository is a mock of ArticleReaderRepository interface.
type MockArticleReaderRepository struct {
	ctrl     *gomock.Controller
	recorder *MockArticleReaderRepositoryMockRecorder
}

// MockArticleReaderRepositoryMockRecorder is the mock recorder for MockArticleReaderRepository.
type MockArticleReaderRepositoryMockRecorder struct {
	mock *MockArticleReaderRepository
}

// NewMockArticleReaderRepository creates a new mock instance.
func NewMockArticleReaderRepository(ctrl *gomock.Controller) *MockArticleReaderRepository {
	mock := &MockArticleReaderRepository{ctrl: ctrl}
	mock.recorder = &MockArticleReaderRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArticleReaderRepository) EXPECT() *MockArticleReaderRepositoryMockRecorder {
	return m.recorder
}

// Save mocks base method.
func (m *MockArticleReaderRepository) Save(ctx context.Context, article domain.Article) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, article)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockArticleReaderRepositoryMockRecorder) Save(ctx, article interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockArticleReaderRepository)(nil).Save), ctx, article)
}

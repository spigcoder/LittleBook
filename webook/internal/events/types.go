package events

const (
	ArticleRead = "article_read"
)

type ReadEvent struct {
	Uid int64
	Aid int64
}

type Consumer interface {
	Start() error
}

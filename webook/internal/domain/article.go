package domain

const (
	ArtStatusUnknown ArtStatus = iota
	ArtStatusPublished
	ArtStatusUnPublished
	ArtStatusPrivate
)

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArtStatus
	CTime   int64
	UTime   int64
}

func (art Article) Abstract() string {
	abs := []rune(art.Content)
	if len(abs) > 100 {
		return string(abs[:100])
	}
	return string(abs)
}

type Author struct {
	Id   int64
	Name string
}

type ArtStatus uint8

func (s ArtStatus) String() string {
	switch s {
	case ArtStatusPublished:
		return "已发布"
	case ArtStatusUnPublished:
		return "未发布"
	case ArtStatusPrivate:
		return "私密"
	default:
		return "未知"
	}
}

func (s ArtStatus) ToUInt8() uint8 {
	return uint8(s)
}

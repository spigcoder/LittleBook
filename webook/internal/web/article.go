package web

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spigcoder/LittleBook/webook/internal/domain"
	"github.com/spigcoder/LittleBook/webook/internal/service"
	"github.com/spigcoder/LittleBook/webook/internal/web/ijwt"
	"net/http"
)

type ArticleHandler struct {
	svc service.ArticleService
}

func NewArticleHandler(svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
	}
}

func (handler *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", handler.Edit)
}

func (handler *ArticleHandler) Edit(c *gin.Context) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Id      int64  `json:"id"`
	}

	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}
	userClaim, ok := c.Get("claims")
	if !ok {
		logrus.Error("claims not found")
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}
	if req.Title == "" || req.Content == "" {
		c.JSON(http.StatusOK, Result{
			Code: http.StatusBadRequest,
			Msg:  "标题或内容为空",
		})
		return
	}
	if len(req.Title) > 1024 {
		c.JSON(http.StatusOK, Result{
			Code: http.StatusBadRequest,
			Msg:  "标题应小于1024",
		})
		logrus.Info("标题长度过长")
	}
	claims, ok := userClaim.(*ijwt.UserClaims)
	//进行内容校验，这里省略
	id, err := handler.svc.Edit(c, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Code: http.StatusInternalServerError,
			Msg:  "系统错误",
		})
		logrus.Error("用户文章发布失败", err)
	}
	c.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "发布成功",
		Data: id,
	})
}

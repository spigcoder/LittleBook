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

type ArticleReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Id      int64  `json:"id"`
}

func (a *ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      a.Id,
		Title:   a.Title,
		Content: a.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}

func NewArticleHandler(svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
	}
}

func (handler *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", handler.Edit)
	g.POST("/publish", handler.Publish)
	g.POST("/withdraw", handler.WithDraw)
}

func (handler *ArticleHandler) WithDraw(c *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}
	//用于验证身份，防止其他人更改我们文章的状态
	userClaim, ok := c.Get("claims")
	if !ok {
		logrus.Error("claims not found")
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}
	claims, ok := userClaim.(*ijwt.UserClaims)
	err := handler.svc.WithDraw(c, req.toDomain(claims.Uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Code: http.StatusInternalServerError,
			Msg:  "系统错误",
		})
		logrus.Error("用户文章对外不可见设置失败", err)
	}
	c.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
	})
}

func (handler *ArticleHandler) Publish(c *gin.Context) {
	var req ArticleReq
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
	id, err := handler.svc.Publish(c, req.toDomain(claims.Uid))
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

func (handler *ArticleHandler) Edit(c *gin.Context) {
	var req ArticleReq
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
	id, err := handler.svc.Save(c, req.toDomain(claims.Uid))
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

package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	regexp "github.com/dlclark/regexp2"
)

type UserHandler struct {
	emailExp *regexp.Regexp
	passExp *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	const (
		EmailRegexp = "^\\w+(-+.\\w+)*@\\w+(-.\\w+)*.\\w+(-.\\w+)*$"
		//8位以上的必须同时包含字母大小写，数字和特殊符号
		PasswordRegexp = "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[!@#$%%^&*()_+{}\\[\\]:;'\",.<>/?\\|~-]).{8,}$"
	)
	//避免每次都要编译
	passRegex := regexp.MustCompile(PasswordRegexp, regexp.None)
	emailRegx := regexp.MustCompile(EmailRegexp, regexp.None)
	
	return &UserHandler{
		emailExp: emailRegx,
		passExp: passRegex,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	s := server.Group("/users")
	s.GET("/profile", u.Profile)
	s.POST("/signup", u.Signup)
	s.POST("/login", u.Login)
	s.POST("/edit", u.Edit)
}

func (u *UserHandler) Profile(c *gin.Context) {

}

func (u *UserHandler) Signup(c *gin.Context) {
	//这个结构体只在当前方法的作用域中有效，出了这个作用域就不可以使用这个结构体了
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var suq SignUpReq
	//Bind方法会根据Content-Type来解析数据
	//如果解析失败会自动返回一个400的错误
	if err := c.Bind(&suq); err != nil {
		return
	}
	if suq.Password != suq.ConfirmPassword {
		c.String(http.StatusBadRequest, "两次密码不一致")
		return
	}
	ok, err := u.emailExp.MatchString(suq.Email)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "服务器问题")
		return
	}
	if !ok {
		c.String(http.StatusBadRequest, "邮箱格式错误")
		return
	}
	ok, err = u.passExp.MatchString(suq.Password)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "服务器问题")
		return
	}
	if !ok {
		c.String(http.StatusBadRequest, "密码必须同时包含字母大小写, 数字和特殊符号, 且不少于8位")
		return
	}
	c.String(200, "注册成功")
}

func (u *UserHandler) Login(c *gin.Context) {

}

func (u *UserHandler) Edit(c *gin.Context) {

}

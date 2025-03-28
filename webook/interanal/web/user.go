package web

import (
	"fmt"
	"net/http"
	"time"
	"unicode/utf8"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/spigcoder/LittleBook/webook/interanal/domain"
	"github.com/spigcoder/LittleBook/webook/interanal/service"
)

var (
	ScretKey = []byte("ZD3oYULPnlBo2wqebduhFQjmrZdaFGaLzayCa8t8HWwxWKbRcGzaNLKkZ31ldeaM")
)

type UserHandler struct {
	emailExp *regexp.Regexp
	passExp  *regexp.Regexp
	svc      *service.UserService
}

type UserClaims struct {
	Uid int64 `json:"uid"`
	jwt.RegisteredClaims
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		EmailRegexp = "^\\w+(-+.\\w+)*@\\w+(-.\\w+)*.\\w+(-.\\w+)*$"
		//8位以上的必须同时包含字母大小写，数字和特殊符号
		PasswordRegexp = "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[!@#$%%^&*()_+{}\\[\\]:;'\",.<>/?\\|~-]).{8,72}$"
	)
	//避免每次都要编译
	passRegex := regexp.MustCompile(PasswordRegexp, regexp.None)
	emailRegx := regexp.MustCompile(EmailRegexp, regexp.None)

	return &UserHandler{
		emailExp: emailRegx,
		passExp:  passRegex,
		svc:      svc,
	}
}

func (handler *UserHandler) RegisterRoutes(server *gin.Engine) {
	s := server.Group("/users")
	s.GET("/profile", handler.Profile)
	s.POST("/signup", handler.Signup)
	// s.POST("/login", handler.Login)
	s.POST("/login", handler.LoginJWT)
	s.POST("/edit", handler.Edit)
}

func (handler *UserHandler) Profile(c *gin.Context) {

}

func (handler *UserHandler) Signup(c *gin.Context) {
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
	ok, err := handler.emailExp.MatchString(suq.Email)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "服务器问题")
		return
	}
	if !ok {
		c.String(http.StatusBadRequest, "邮箱格式错误")
		return
	}
	ok, err = handler.passExp.MatchString(suq.Password)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "服务器问题")
		return
	}
	if !ok {
		c.String(http.StatusBadRequest, "密码必须同时包含字母大小写, 数字和特殊符号, 且不少于8位")
		return
	}

	//调用服务接口
	err = handler.svc.SignUp(c, domain.User{
		Email:    suq.Email,
		Password: suq.Password})
	if err == service.ErrDuplicateEmail {
		c.String(http.StatusBadRequest, "邮箱冲突")
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器问题")
		return
	}
	c.String(200, "注册成功")
}
func (handler *UserHandler) LoginJWT(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Id       int64  `json:"id"`
	}

	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return
	}
	u, err := handler.svc.Login(c, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrInvalidUserOrPassword {
		c.String(http.StatusBadRequest, "用户名或密码错误")
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器问题")
		return
	}
	//设置JWT
	userClaims := UserClaims{
		Uid: u.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaims)
	tokenStr, err := token.SignedString([]byte(ScretKey))
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器问题")
		return
	}
	c.Header("x-jwt-token", tokenStr)
	c.String(200, "登录成功")
}

func (handler *UserHandler) Login(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Id       int64  `json:"id"`
	}

	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return
	}
	u, err := handler.svc.Login(c, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrInvalidUserOrPassword {
		c.String(http.StatusBadRequest, "用户名或密码错误")
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器问题")
		return
	}
	//设置session
	see := sessions.Default(c)
	see.Set("userId", u.Id)
	see.Options(sessions.Options{
		MaxAge: 60 * 30,
	})
	see.Save()
	c.String(200, "登录成功")
}

func (handler *UserHandler) Edit(c *gin.Context) {
	type EditReq struct {
		UserName string `json:"userName"`
		Birthday string `json:"birthday"`
		Intro    string `json:"intro"`
	}
	var req EditReq
	if err := c.Bind(&req); err != nil {
		return
	}
	//校验
	if utf8.RuneCountInString(req.UserName) > 32 {
		c.String(http.StatusBadRequest, "用户名长度不能超过32")
		return
	}
	if utf8.RuneCountInString(req.Intro) > 256 {
		c.String(http.StatusBadRequest, "简介长度不能超过256")
		return
	}
	userClaim, ok := c.Get("claims")
	if !ok {
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}
	claims, ok := userClaim.(*UserClaims)
	if !ok {
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}
	id := claims.Uid
	err := handler.svc.Edit(c, domain.User{
		UserName: req.UserName,
		Birthday: req.Birthday,
		Intro:    req.Intro,
		Id:       id,
	})

	if err != nil {
		c.String(http.StatusInternalServerError, "服务器问题")
		return
	}
	c.String(200, "修改成功")
}

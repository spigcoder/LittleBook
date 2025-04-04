package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spigcoder/LittleBook/webook/internal/service"
	svcmocks "github.com/spigcoder/LittleBook/webook/internal/service/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_Signup(t *testing.T) {
	testClass := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				us := svcmocks.NewMockUserService(ctrl)
				us.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				return us
			},
			reqBody: `{
			"email": "2240434341@qq.com",
			"password": "Ta.123456",
			"confirmPassword": "Ta.123456"
			}`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "绑定失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				us := svcmocks.NewMockUserService(ctrl)
				return us
			},
			reqBody: `{
			"email": "2240434341@qq.com",
			"password": "Ta.123456"
			"confirmPassword": "Ta.123456"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: "",
		},
		{
			name: "密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				us := svcmocks.NewMockUserService(ctrl)
				return us
			},
			reqBody: `{
			"email": "2240434341@qq.com",
			"password": "Ta.123456",
			"confirmPassword": "Ta.1234561"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: "两次密码不一致",
		},
		{
			name: "邮箱格式错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				us := svcmocks.NewMockUserService(ctrl)
				return us
			},
			reqBody: `{
			"email": "22404.34341@qq.com",
			"password": "Ta.123456",
			"confirmPassword": "Ta.123456"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: "邮箱格式错误",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				us := svcmocks.NewMockUserService(ctrl)
				us.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(service.ErrDuplicateEmail)
				return us
			},
			reqBody: `{
			"email": "2240434341@qq.com",
			"password": "Ta.123456",
			"confirmPassword": "Ta.123456"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: "邮箱冲突",
		},
	}
	for _, tc := range testClass {
		t.Run(tc.name, func(t *testing.T) {
			server := gin.Default()
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			h := NewUserHandler(tc.mock(ctl), nil)
			h.RegisterRoutes(server)
			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			t.Log(resp)
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}

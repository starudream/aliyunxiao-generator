package api

import (
	"fmt"
	"net/http"

	"github.com/starudream/go-lib/core/v2/config"
	"github.com/starudream/go-lib/resty/v2"
)

const Addr = "https://devops.aliyun.com/projex/api"

type BaseResp[T any] struct {
	Success *bool  `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`

	Code     *int   `json:"code"`
	ErrorMsg string `json:"errorMsg,omitempty"`

	Result T `json:"result,omitempty"`
}

func (t *BaseResp[T]) IsSuccess() bool {
	return t != nil && ((t.Success != nil && *t.Success) || (t != nil && *t.Code == 200))
}

func (t *BaseResp[T]) String() string {
	if t != nil && t.Success != nil {
		return fmt.Sprintf("error: %s, msg: %s", t.Error, t.Message)
	} else if t != nil && t.Code != nil {
		return fmt.Sprintf("code: %d, msg: %s", *t.Code, t.ErrorMsg)
	}
	return "nil"
}

func R() *resty.Request {
	return resty.R().
		SetHeader("User-Agent", resty.UAWindowsChrome).
		SetHeader("Accept-Encoding", "gzip").
		SetCookie(&http.Cookie{Name: "login_aliyunid_ticket", Value: config.Get("ticket").String()})
}

func Exec[T any](r *resty.Request, method, path string) (t T, _ error) {
	res, err := resty.ParseResp[*BaseResp[any], *BaseResp[T]](
		r.SetError(&BaseResp[any]{}).SetResult(&BaseResp[T]{}).Execute(method, Addr+path),
	)
	if err != nil {
		return t, fmt.Errorf("[api] %w", err)
	}
	return res.Result, nil
}

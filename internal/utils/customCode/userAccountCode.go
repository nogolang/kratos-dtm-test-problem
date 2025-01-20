package customCode

import (
	"github.com/go-kratos/kratos/v2/errors"
	"net/http"
)

var (
	//使用http的状态码，kratos会自动转换
	//统一返回200
	UserAccountNotFount  = errors.New(http.StatusNotFound, "UserAccountNotFount", "不存在")
	UserAccountDeleteIds = errors.New(http.StatusBadRequest, "UserAccountIdsNeed", "需要指定要删除的id")
)

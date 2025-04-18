package handler

import (
	"net/http"

	"fogate/internal/logic"
	"fogate/internal/svc"
	"fogate/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FogateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewFogateLogic(r.Context(), svcCtx)
		resp, err := l.Fogate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

package handler

import (
	"net/http"

	"fogate/internal/svc"
	"fogate/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ToolsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 设置默认值
		if req.Page < 1 {
			req.Page = 1
		}
		if req.PageSize < 1 || req.PageSize > 100 {
			req.PageSize = 20
		}

		// 使用读锁获取工具信息
		svcCtx.ToolsCache.RLock()
		defer svcCtx.ToolsCache.RUnlock()

		// 计算分页
		start := (req.Page - 1) * req.PageSize
		end := start + req.PageSize
		if start >= len(svcCtx.ToolsCache.Tools) {
			// 如果起始位置超出范围，返回空列表
			httpx.OkJson(w, types.ToolsResponse{
				Total: svcCtx.ToolsCache.Total,
				Tools: []types.SimpleToolInfo{},
			})
			return
		}
		if end > len(svcCtx.ToolsCache.Tools) {
			end = len(svcCtx.ToolsCache.Tools)
		}

		// 返回分页后的工具信息
		httpx.OkJson(w, types.ToolsResponse{
			Total: svcCtx.ToolsCache.Total,
			Tools: svcCtx.ToolsCache.Tools[start:end],
		})
	}
}

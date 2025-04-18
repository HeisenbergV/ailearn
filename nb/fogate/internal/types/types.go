package types

type Request struct {
	Name string `path:"name,options=you|me"`
}

type Response struct {
	Message string `json:"message"`
}

// 用于解析原始 JSON 的结构体
type Tool struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Capability struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ToolInfo struct {
	Tool         Tool         `json:"tool"`
	Capabilities []Capability `json:"capabilities"`
}

// 简化后的响应结构体
type SimpleCapability struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SimpleToolInfo struct {
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Capabilities []SimpleCapability `json:"capabilities"`
}

// 分页请求
type PageRequest struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"pageSize,default=20"`
}

// 分页响应
type ToolsResponse struct {
	Total int              `json:"total"`
	Tools []SimpleToolInfo `json:"tools"`
}

package logic

import (
	"context"

	"fogate/internal/svc"
	"fogate/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FogateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFogateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FogateLogic {
	return &FogateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FogateLogic) Fogate(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}

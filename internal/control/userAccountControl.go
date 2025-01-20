package controller

import (
	"context"
	"database/sql"
	_ "github.com/dtm-labs/driver-kratos"
	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"kraots-xa/configs"
	"kraots-xa/internal/service"
	"kraots-xa/internal/utils/dtmUtils"
	"kraots-xa/proto/pb/userAccountPb"
)

type UserAccountControl struct {
	userAccountPb.UnimplementedUserAccountServer
	UserAccountSvc *service.UserAccountSvcImp
	AllConfig      *configs.AllConfig
	Logger         *zap.Logger
}

func NewUserAccountControl(grpcServer *grpc.Server,
	httpServer *http.Server,
	logger *zap.Logger,
	allConfig *configs.AllConfig,
	userAccountSvc *service.UserAccountSvcImp) *UserAccountControl {

	u := new(UserAccountControl)
	u.Logger = logger
	u.UserAccountSvc = userAccountSvc
	u.AllConfig = allConfig

	//注册grpc server
	userAccountPb.RegisterUserAccountServer(grpcServer, u)

	//注册http server
	userAccountPb.RegisterUserAccountHTTPServer(httpServer, u)

	return u
}

// 此时kratos会从etcd里获取到dtm要请求的服务的地址
// @TODO 这个有问题，但是dtm是可以和这个服务通信的，并且已经调用了分支事务，但是提交的时候dtm本身却报错
var baseUrl = "discovery:///dtm-test-service"

// 而dtm的本身的地址，也可以通过服务发现获取，前提dtm已经注册到etcd中了
var DtmServer = "discovery:///dtmservice" //这个没问题

// grpc连接则用这种，无需加任何前缀
//var baseUrl = "192.168.80.1:5008"

func (receiver *UserAccountControl) UpdateAccount(ctx context.Context, request *userAccountPb.UserAccountTransRequest) (*emptypb.Empty, error) {
	//调用各个分支事务

	//要转的金额
	amount := request.Amount

	//从用户1转到用户2
	out := &userAccountPb.UserAccountUpdateRequest{
		Uid:            request.Uid,
		Amount:         amount,
		TransOutResult: request.TransOutResult,
	}

	in := &userAccountPb.UserAccountUpdateRequest{
		Uid:           request.Tid,
		Amount:        amount,
		TransInResult: request.TransInResult,
	}

	gid := dtmgrpc.MustGenGid(DtmServer)
	err := dtmgrpc.XaGlobalTransaction(DtmServer, gid, func(xa *dtmgrpc.XaGrpc) error {
		//resp里的body取决于接口返回什么，它是grpc的返回值
		//但是它的返回值是空，所以这里用空即可
		r := &emptypb.Empty{}
		err := xa.CallBranch(out, baseUrl+"/UserAccount/TransOutXa", r)
		if err != nil {
			receiver.Logger.Error("TransOutXa出错", zap.Error(err))
			return err
		}

		err = xa.CallBranch(in, baseUrl+"/UserAccount/TransInXa", r)
		if err != nil {
			receiver.Logger.Error("TransInXa出错", zap.Error(err))
			return err
		}
		return nil
	})
	receiver.Logger.Info("gid==", zap.String("", gid))
	if err != nil {
		receiver.Logger.Error("全局事务出错", zap.Error(err))
	}

	return &emptypb.Empty{}, err
}

func (receiver *UserAccountControl) TransOutXa(ctx context.Context, request *userAccountPb.UserAccountUpdateRequest) (*emptypb.Empty, error) {
	//注册本地分支事务RM
	err := dtmgrpc.XaLocalTransaction(ctx, receiver.AllConfig.DtmConf, func(db *sql.DB, xa *dtmgrpc.XaGrpc) error {

		//转出操作
		request.Amount = -request.Amount
		//这里用自己的db对象，但是需要封装一下dtm提供的db
		//@TODO 这里是没问题的，用dtm提供的db也一样
		gormDb := dtmUtils.GetGormDbFromDtmConn(db, receiver.Logger)
		return receiver.UserAccountSvc.UpdateAccount(gormDb, request)
	})

	if err != nil {
		//里面会处理上面dtm返回的错误码，填充message
		err = dtmgrpc.GrpcError2DtmError(err)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (receiver *UserAccountControl) TransInXa(ctx context.Context, request *userAccountPb.UserAccountUpdateRequest) (*emptypb.Empty, error) {
	//注册本地分支事务RM
	err := dtmgrpc.XaLocalTransaction(ctx, receiver.AllConfig.DtmConf, func(db *sql.DB, xa *dtmgrpc.XaGrpc) error {
		//转入操作，Amount是正数
		request.Amount = request.Amount

		//这里用自己的db对象，但是需要封装一下dtm提供的db
		gormDb := dtmUtils.GetGormDbFromDtmConn(db, receiver.Logger)
		return receiver.UserAccountSvc.UpdateAccount(gormDb, request)
	})

	//里面会处理上面dtm返回的错误码，填充message
	if err != nil {
		//里面会处理上面dtm返回的错误码，填充message
		err = dtmgrpc.GrpcError2DtmError(err)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

package test

import (
	"context"
	"kraots-xa/proto/pb/userAccountPb"
	"testing"
)

func Test_TransMoney(t *testing.T) {
	grpcClient := GetGrpcClient()

	client := userAccountPb.NewUserAccountClient(grpcClient)

	var request userAccountPb.UserAccountTransRequest
	//用户1向用户2转账10元
	request.Uid = 1
	request.Tid = 2
	request.Amount = 10
	request.TransOutResult = ""
	request.TransInResult = ""
	//request.TransInResult = dtmcli.ResultFailure

	ctx := context.Background()
	_, err := client.UpdateAccount(ctx, &request)
	if err != nil {
		t.Error(err)
	}
}

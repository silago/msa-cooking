package server

import (
	"context"
	"cooking-users/proto"
	"cooking-users/user_data_providers"
)

type Server struct {
	provider user_data_providers.UserDataProvider
}

func (server Server) GetUserResource(c context.Context, r *cooking_users.ResourceRequest) (*cooking_users.ResourceResponse, error) {
	resource :=server.provider.GetOne(r.UserId, r.Type);
	response := cooking_users.ResourceResponse{ Value: resource}
	return &response, nil
}

func (server Server) DrawUserResource(ctx context.Context, request *cooking_users.ChangeResourceRequest) (*cooking_users.ChangeResourceResponse, error) {
	result, err:= server.provider.IncrementOne(request.UserId,request.Type, -request.Count);
	msg:=""
	if err == nil {
		msg=result
	} else {
		msg = err.Error()
	}
	return &cooking_users.ChangeResourceResponse{
		Result:err==nil,
		Msg:msg,
	}, err
}

func (server Server) AddUserResource(ctx context.Context, request *cooking_users.ChangeResourceRequest) (*cooking_users.ChangeResourceResponse, error) {
	result, err:= server.provider.IncrementOne(request.UserId,request.Type, request.Count)
	msg:=""
	if err == nil {
		msg=result
	} else {
		msg = err.Error()
	}
	return &cooking_users.ChangeResourceResponse{
		Result:err==nil,
		Msg:msg,
	}, err
}




package main

import (
	"context"
	"github.com/silago/msa-cooking/cooking-users/proto"
	"github.com/silago/msa-cooking/cooking-users/user_data_providers"
)


type UserResourceServer struct {
	dataProvider user_data_providers.UserDataProvider
}

func (u *UserResourceServer) GetUserResource(ctx context.Context, request *cooking_users.ResourceRequest) (*cooking_users.ResourceResponse, error) {
	data:=u.dataProvider.GetOne(request.UserId,request.Type)
	response:=cooking_users.ResourceResponse{}
	response.Value=data
	return &response, nil
}

func (u *UserResourceServer) DrawUserResource(ctx context.Context, request *cooking_users.ChangeResourceRequest) (*cooking_users.ChangeResourceResponse, error) {
	result, err:= u.dataProvider.IncrementOne(request.UserId, request.Type, -request.Count)
	response:=cooking_users.ChangeResourceResponse{}
	response.Result = err==nil
	response.Msg=result
	return &response,err
}

func (u *UserResourceServer) UpdateUserResources(ctx context.Context, request *cooking_users.ResourceRequest) (*cooking_users.ResourceResponse, error) {
	panic("implement me")
	//	u.dataProvider.SetMany(request.UserId)
}

func (u *UserResourceServer) AddUserResource(ctx context.Context, request *cooking_users.ChangeResourceRequest) (*cooking_users.ChangeResourceResponse, error) {
	result, err:= u.dataProvider.IncrementOne(request.UserId, request.Type, request.Count)
	response:=cooking_users.ChangeResourceResponse{}
	response.Result = err==nil
	response.Msg=result
	return &response,err
}

func (u *UserResourceServer) AddUserResources(ctx context.Context, request *cooking_users.ChangeResourcesRequest) (*cooking_users.ResourcesResponse, error) {
	updateData:= make(map[string]interface{})
	for _, resource :=  range request.Resources {
		updateData[resource.Type]= resource.Value
	}
	if e := u.dataProvider.IncrementMany(request.UserId, updateData); e !=nil {
		return nil, e
	} else {
		response:=cooking_users.ResourcesResponse{}
		response.Resources = make([]*cooking_users.Resource, len(request.Resources))
		data:=u.dataProvider.GetAll(request.UserId)
		for index,resource:=range request.Resources {
			response.Resources[index]=&cooking_users.Resource{Type:resource.Type,Value:data[resource.Type]}
		}
		return &response, nil
	}
}

package main

import (
	"context"
	"cooking-users/modules/locations"
	"cooking-users/proto"
	"cooking-users/user_data_providers"
	. "cooking-users/user_service"
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/nats-io/go-nats"
	"log"
	"net/http"
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

func main() {
	dbConnection := redis.NewClient(&redis.Options{
		Addr: ENV("REDIS_HOST"), //"localhost:6379",
		DB: 0, // use default DB
	})
	natsUrl := ENV("NATS_HOST")
	var natsConnection *nats.Conn
	var natsService *nats.EncodedConn
	var err error

	if natsConnection, err = nats.Connect(natsUrl); err != nil {
		log.Panicf("nats host is not available: {%s}", err)
	} else if natsService, err = nats.NewEncodedConn(natsConnection, nats.JSON_ENCODER); err != nil {
		log.Panicf("nats connecton is not available: {%s}", err)
	}
	userProvider := user_data_providers.NewRedisUserProvider(dbConnection) //NewRedisUserProvider(dbConnection)//&UserProvider{Client: dbConnection}
	if _, err := natsService.Subscribe(USER_INFO_CHANNEL, func(subj string, message *NatsMessage) {
		var data map[string]interface{}
		json.Unmarshal([]byte(message.Message.(string)), &data)
		switch message.Type {
		case NATS_SET_MESSAGE:
			userProvider.SetMany(message.UserId, data)
		case NATS_APPEND_MESSAGE:
			userProvider.IncrementMany(message.UserId, data)
		}
	}); err != nil {
		log.Println("nats subscription error", err)
	}

	userService:=NewUserService(userProvider, natsService);

	http.HandleFunc("/update/", userService.Update())
	http.HandleFunc("/signin/",userService.Auth())
	http.HandleFunc("/auth/",  userService.Auth())
	http.HandleFunc("/drop/",  userService.Drop())
	http.HandleFunc("/healthcheck/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("alive"))
	})

	if locationHandler, moduleError:= locations.NewLocationUpdateHandler("assets/resources/locations.xml", userProvider); moduleError == nil {
		log.Printf("prepare to handle /game/location/unlock/")
		http.HandleFunc("/game/location/unlock/",func(w http.ResponseWriter, r *http.Request) {
			if userId:= r.Header.Get(USER_ID_HEADER); userId == "" {
				w.WriteHeader(401)
				log.Printf("401 ")
				return
			} else {
				result:= locationHandler(userId, r)
				if _, e:= w.Write([]byte(result)); e!=nil {
					log.Printf("response write error {%s} " , e.Error())
				}
			}
		})
		log.Printf("Handling /game/location/unlock/")
	} else {
		log.Fatalf("error on location module start {%s}", moduleError.Error())
	}

	cooking_users.NewUserResourcesServer( &UserResourceServer{}, nil )

	log.Println("listening at", ENV("HOST"))
	log.Fatal(http.ListenAndServe(ENV("HOST"), nil))
}





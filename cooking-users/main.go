package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/nats-io/go-nats"
	"github.com/silago/msa-cooking/cooking-users/modules/locations"
	"github.com/silago/msa-cooking/cooking-users/proto"
	"github.com/silago/msa-cooking/cooking-users/user_data_providers"
	. "github.com/silago/msa-cooking/cooking-users/user_service"
	"log"
	"net/http"
)


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


	if locationsModule, moduleError:=locations.NewLocationsModule("assets/resources/locations.xml", userProvider); moduleError!=nil {
		log.Fatalf("error on location module start {%s}", moduleError.Error())
	} else {
		http.HandleFunc("/game/location/unlock/",func(w http.ResponseWriter, r *http.Request) {
			if userId:= r.Header.Get(USER_ID_HEADER); userId != "" {
				_, _= w.Write([]byte(locationsModule.UnlockHandler(userId, r)))
			}
		})

		http.HandleFunc("/game/location/upgrade/",func(w http.ResponseWriter, r *http.Request) {
			if userId:= r.Header.Get(USER_ID_HEADER); userId != "" {
				_, _= w.Write([]byte(locationsModule.UpgradeHandler(userId, r)))
			}
		})
	}

	twirpHandler:=cooking_users.NewUserResourcesServer( &UserResourceServer{dataProvider:userProvider}, nil )
	mux:=http.NewServeMux()
	mux.Handle(cooking_users.UserResourcesPathPrefix, twirpHandler)
	go func() {
		log.Printf("Twirp listening at %s", ENV("TWIRP_ADDR"))
		log.Fatal(http.ListenAndServe(ENV("TWIRP_ADDR"), mux))
	}()

	log.Println("listening at", ENV("HOST"))
	log.Fatal(http.ListenAndServe(ENV("HOST"), nil))
}






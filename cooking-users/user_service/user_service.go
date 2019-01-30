package user_service

import (
	//socialAuth "social-auth-go/api"
	socialAuth "github.com/silago/social-api"
	"encoding/json"
	"github.com/nats-io/go-nats"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	//. "../user_data_providers"
	providers "cooking-users/user_data_providers"
)

type NatsMessageType int
type NatsMessage struct {
	UserId  string
	Type    NatsMessageType
	Message interface{}
}


const USER_ID_HEADER = "X-User-Id"
const USER_FRIENDS_CHANNEL = "user.friends.change"
const USER_INFO_CHANNEL = "user.params.change"
const NATS_MESSAGE_DEFAULT = 0
const NATS_SET_MESSAGE = 1
const NATS_APPEND_MESSAGE = 2

func ENV(name string) string {
	result := ""
	if s, ok := os.LookupEnv(name); ok {
		result = s
	} else {
		log.Fatal("Could not get env var " + name)
	}
	return result
}

type AuthRequest struct {
	Platform   	  string `json:"platform,omitempty"`
	SessionKey    string `json:"session_key,omitempty"`
	SessionSecret string `json:"session_secret,omitempty"`
	SessionToken string `json:"session_token,omitempty"`
}

type UserService struct {
	userProvider providers.UserDataProvider;
	nats *nats.EncodedConn;
}

func NewUserService(provider providers.UserDataProvider, nats *nats.EncodedConn) *UserService {
	service:= UserService{provider, nats}
	return &service
}

func GetAuthenticator(request AuthRequest) socialAuth.AuthProvider {
	var   authenticator socialAuth.AuthProvider
	switch request.Platform {
	case "mock":
		authenticator = socialAuth.NewMockAuthProvider();
		break
	case "ok", "":
		authenticator = socialAuth.NewOkAuthProvider(ENV("APP_ID"),request.SessionKey,request.SessionSecret)
		break
	case "vk":
		authenticator = socialAuth.NewVkAuthProvider(request.SessionToken)
		break

	}
	return authenticator;
}


func (service *UserService) Auth() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := AuthRequest{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&request); err != nil {
			log.Println("[ERROR] error request decode: ", err.Error())
			return
		}

		var authenticator socialAuth.AuthProvider=GetAuthenticator(request)
		platformUser, err := authenticator.Auth()
		if err != nil {
			log.Printf("[ERROR] error: {%s}", err.Error())
			_, _ = w.Write([]byte("{\"error\":\"ok_auth failed \"}"))
			return
		}

		if platformUser.Uid == "" {
			_, _ = w.Write([]byte("{\"error\":\"ok_auth failed user id is empty\"}"))
			return
		}
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		if err == nil {
			if user, err := service.userProvider.GetOrCreate(platformUser.Uid, map[string]interface{}{"username": platformUser.FirstName + " " + platformUser.LastName}); err == nil {
				rating := user.Params["user_rating"]
				infoMessage, _ := json.Marshal(map[string]interface{}{"user_id": platformUser.Uid, "username": platformUser.FirstName + " " + platformUser.LastName, "pic": platformUser.PicBase, "rating": strconv.Itoa(rating.(int))})
				natsErr := service.nats.Publish(USER_INFO_CHANNEL, NatsMessage{UserId: platformUser.Uid, Message: string(infoMessage)})
				paramsMessage, _ := json.Marshal(user.Params)
				log.Printf("params message: %v", user.Params)
				_ = service.nats.Publish(USER_INFO_CHANNEL, NatsMessage{UserId: platformUser.Uid, Type: NATS_SET_MESSAGE, Message: string(paramsMessage)})

				if natsErr != nil {
					log.Println("[ERROR] nats Publish error", natsErr.Error())
				} else {
					log.Printf("Nats data published", infoMessage)
				}
				if response, err := json.Marshal(user.Params); err == nil {
					w.Write([]byte(response))
				}
			} else {
				log.Printf("[ERROR] User init error: ", err.Error())
				w.Write([]byte(err.Error()))
			}
			if friendsData, err := authenticator.Friends(); err != nil {
				log.Printf("[ERROR] error get friends:: ", friendsData, err)
			} else if friendsMessage, err := json.Marshal(friendsData); err != nil {
			} else {
				service.nats.Publish(USER_FRIENDS_CHANNEL, NatsMessage{UserId: platformUser.Uid, Message: string(friendsMessage)})
			}
		} else {
			log.Println("ok_auth error ", err)
		}
	}
}

func (service *UserService) NatsAuth() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := AuthRequest{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&request); err != nil {
			log.Println("[ERROR] error request decode: ", err.Error())
			return
		}

		//authenticator := socialAuth.Api{
		//	AppId: ENV("APP_ID"),
		//}
		authenticator:= GetAuthenticator(request)

		platform_user, err := authenticator.Auth()
		if err != nil {
			log.Printf("[ERROR] error: {%s}", err.Error())
			w.Write([]byte("{\"error\":\"ok_auth failed\"}"))
			return
		}

		if platform_user.Uid == "" {
			w.Write([]byte("{\"error\":\"ok_auth failed\"}"))
			return
		}
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		if err == nil {
			if user, err := service.userProvider.GetOrCreate(platform_user.Uid, map[string]interface{}{"username": platform_user.FirstName + " " + platform_user.LastName}); err == nil {
				rating := user.Params["user_rating"]
				if rating == nil {
					rating = ""
				}

				infoMessage, _ := json.Marshal(map[string]string{"user_id": platform_user.Uid, "username": platform_user.FirstName + " " + platform_user.LastName, "pic": platform_user.PicBase, "rating": rating.(string)})
				nats_err := service.nats.Publish(USER_INFO_CHANNEL, NatsMessage{UserId: platform_user.Uid, Message: string(infoMessage)})
				paramsMessage, _ := json.Marshal(user.Params)
				log.Printf("params message:", user.Params)
				service.nats.Publish(USER_INFO_CHANNEL, NatsMessage{UserId: platform_user.Uid, Type: NATS_SET_MESSAGE, Message: string(paramsMessage)})

				if nats_err != nil {
					log.Println("[ERROR] nats Publish error", nats_err.Error())
				} else {
					log.Printf("Nats data published", infoMessage)
				}
				if response, err := json.Marshal(user.Params); err == nil {
					w.Write([]byte(response))
				}
			} else {
				log.Printf("[ERROR] User init error: ", err.Error())
				w.Write([]byte(err.Error()))
			}
			if friendsData, err := authenticator.Friends(); err != nil {
				log.Printf("[ERROR] error get friends:: ", friendsData, err)
			} else if friendsMessage, err := json.Marshal(friendsData); err != nil {
			} else {
				service.nats.Publish(USER_FRIENDS_CHANNEL, NatsMessage{UserId: platform_user.Uid, Message: string(friendsMessage)})
			}
		} else {
			log.Println("ok_auth error ", err)
		}
	}
}



func (service *UserService) Update() func(w http.ResponseWriter, r *http.Request) {
	userProvider:=service.userProvider;
	nats:=service.nats;
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err    error
			user   providers.User
			userId string
		)
		if userId = r.Header.Get(USER_ID_HEADER); userId == "" {
			w.WriteHeader(401)
			return
		}
		body, _ := ioutil.ReadAll(r.Body)

		data := make(map[string]interface{})
		err = json.Unmarshal(body, &data)
		if user, err = userProvider.GetOrCreate(userId, nil); err != nil {
			log.Println("User retrievveing failure: ", err.Error())
		} else {
			err = userProvider.SetMany(user.UserId, data)
			if rating := calcRatingByParams(data, user); rating != 0 {
				if newRating, _ := userProvider.IncrementOne(user.UserId, "user_rating", float64(rating)); newRating != "0" {
					data["user_rating"] = newRating
				}
			}

			msg, _ := json.Marshal(data)
			if publishErr:= nats.Publish(USER_INFO_CHANNEL, NatsMessage{UserId: user.UserId, Type: NATS_SET_MESSAGE, Message: string(msg)}); publishErr!=nil {
				log.Println("Nats publish error %s", publishErr.Error())
			}
		}
	}
}

func calcRatingByParams(params map[string]interface{}, user providers.User) int {
	ratingKeys := []string{
		"bakery_level",
		"bakery_level_1",
		"bakery_level_2",
		"bakery_level_3",
		"bakery_level_4",
		"bakery_level_5",
		"bakery_level_6",
		"bakery_level_7",
		"bakery_level_8",
		"bakery_level_9",
		"bakery_level_10",
		"bakery_level_11",
		"bakery_level_12",
		"bakery_level_13",
		"bakery_level_14",
		"bakery_level_15",
		"bakery_level_16",
		"bakery_level_17",
		"bakery_level_18",
	}
	var rating int
	for _, key := range ratingKeys {
		if params[key] != nil {
			//fmt.Println(">>", params[key])
			newVal := int(params[key].(float64))
			oldVal := int(0)
			if user.Params[key] != nil {
				oldVal, _ = strconv.Atoi(user.Params[key].(string))
			}
			rating += (newVal - oldVal)
		}
	}
	return int(rating)
}


func (service *UserService) Drop() func(w http.ResponseWriter, r *http.Request) {
	userProvider:=service.userProvider;
	return func(w http.ResponseWriter, r *http.Request) {
		if token:= r.Header.Get("X-Adm-Token"); token != ENV("ADMINISTRATION_TOKEN") {
			w.WriteHeader(401)
			return
		} else {
			body, _ := ioutil.ReadAll(r.Body)
			data := make(map[string]interface{})
			if err:=json.Unmarshal(body, &data); err==nil {
				userProvider.DeleteAll(data["id"].(string))
			}
		}
		w.Write([]byte("User dropped"))
	}
}


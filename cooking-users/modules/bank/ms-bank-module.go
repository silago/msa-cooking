package bank
/*
import (
	"database/sql"
	"github.com/go-redis/redis"
	"github.com/nats-io/go-nats"
	"crypto/md5"
	"encoding/json"
	"sort"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
)

const REDIS_STATS_KEY = "U:"

///struct Items map[string]float64;
var ProductList = map[string][]Product{
	"coins-300":     {{Type: "coins", Count: 300}, {Type: "candies", Count: 5}},
	"coins-1500":    {{Type: "coins", Count: 1500}, {Type: "candies", Count: 5}},
	"coins-3000":    {{Type: "coins", Count: 3000}, {Type: "candies", Count: 5}},
	"coins-8000":    {{Type: "coins", Count: 8000}, {Type: "candies", Count: 5}},
	"crystals-50":   {{Type: "crystals", Count: 50}, {Type: "candies", Count: 5}},
	"crystals-150":  {{Type: "crystals", Count: 150}, {Type: "candies", Count: 5}},
	"crystals-700":  {{Type: "crystals", Count: 700}, {Type: "candies", Count: 5}},
	"crystals-1500": {{Type: "crystals", Count: 1500}, {Type: "candies", Count: 5}},
}

type NatsMessageType int

const NATS_MESSAGE_DEFAULT = 0
const NATS_SET_MESSAGE = 1
const NATS_APPEND_MESSAGE = 2

type NatsMessage struct {
	UserId  string
	Type    NatsMessageType
	Message interface{}
}

const USER_ID_HEADER = "X-User-Id"
const USER_FRIENDS_CHANNEL = "user.friends.change"

const USER_INFO_CHANNEL = "user.params.change"

type Product struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

type Transaction struct {
	TransactionId string
	UserId        string
	ItemId        string //varchar(33) NOT NULL,
	Complete      bool
}

// функция рассчитывает подпись для пришедшего запроса
// подробнее про алгоритм расчета подписи можно посмотреть в документации (https://apiok.ru/dev/methods/)

func calcSignature(r *http.Request) string {
	params := make(map[string]string)
	for k, v := range r.URL.Query() {
		if k == "sig" {
			continue
		} else {
			params[k] = v[0]
		}
	}

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var str string
	for _, k := range keys {
		str = str + k + "=" + params[k]
	}
	str = str + ENV("APP_SECRET")

	data := []byte(str)
	result := fmt.Sprintf("%x", md5.Sum(data))
	return result
}

func GetAllIncomplete(db *sql.DB, UserId string) ([]Product, error) {
	var err error
	items := make([]Product, 0)
	if result, err := db.Query("select item_id from payments where complete = false and user_id = ?", UserId); err == nil {
		defer result.Close()
		for result.Next() {
			var productId string
			result.Scan(&productId)
			products := ProductList[productId]
			items = append(items, products...)
		}
	} else {
		return nil, err
	}
	return items, err
}

func CompleteAll(db *sql.DB, UserId string) error {
	var err error
	if result, err := db.Query("update payments set complete = TRUE where user_id = ?", UserId); err != nil {
		defer result.Close()
	}
	return err
}

func (t *Transaction) Create(db *sql.DB) error {
	var err error
	if stm, err := db.Prepare("insert into payments values (?, ?, ?, ?, NOW() )"); err == nil {
		_, err = stm.Exec(t.TransactionId, t.UserId, t.ItemId, t.Complete)
		return err
	} else {
		return err
	}
	return err
}

func ENV(name string) string {
	if s, ok := os.LookupEnv(name); ok {
		return s
	}
	return ""
}

func successResponse() string {
	return `
    <?xml version="1.0" encoding="UTF-8"?>
    <callbacks_payment_response xmlns="http://api.forticom.com/1.0/">\n
        true\n
    </callbacks_payment_response>`
}

func errorResponse(code string, messageCode string, messageText string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
        <ns2:error_response xmlns:ns2='http://api.forticom.com/1.0/'>
            <error_code>%s</error_code>
            <error_msg>%s : %s</error_msg>
        </ns2:error_response>`, code, messageCode, messageText)
}

func get(r *http.Request, key string) string {
	keys, ok := r.URL.Query()[key]
	if !ok || len(keys[0]) < 1 {
		return ""
	}
	return keys[0]
}

func receive(db *redis.Client, sql *sql.DB, nats *nats.EncodedConn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var UserId string
		if UserId = r.Header.Get(USER_ID_HEADER); UserId == "" {
			//response.Code = 401
			w.WriteHeader(401)
			return
		}
		products, err := GetAllIncomplete(sql, UserId)
		if err != nil {
			log.Printf("Error getting payments: ", err.Error())
		}
		result := make(map[string]int64)
		for _, product := range products {
			result["congrat_"+product.Type] += product.Count
			newValue, _ := db.HGet(REDIS_STATS_KEY+UserId, product.Type).Int64()
			newValue += product.Count
			result[product.Type] += newValue
			params := map[string]int64{product.Type: product.Count}
			msg, _ := json.Marshal(params)
			nats.Publish(USER_INFO_CHANNEL, NatsMessage{UserId: UserId, Type: NATS_APPEND_MESSAGE, Message: string(msg)})
			CompleteAll(sql, UserId)
		}
		jsonResult, _ := json.Marshal(result)
		w.Write(jsonResult)
	}
}

func callback(db *redis.Client, sql *sql.DB, nats *nats.EncodedConn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var response string
		defer func() {
			w.Write([]byte(response))
		}()

		productId := get(r, "product_code")
		amount := get(r, "amount")
		sig := get(r, "sig")
		uid := get(r, "uid")
		transactionId := get(r, "transaction_id")
		if productId == "" || amount == "" || uid == "" || sig == "" || uid == "" {
			response = errorResponse("1001", "CALLBACK_INVALID_PAYMENT", fmt.Sprintf("Wrong callback params", productId, amount, sig, uid, transactionId))
			return
		}

		if sig != calcSignature(r) {
			response = errorResponse("1001", "CALLBACK_INVALID_PAYMENT", "invalid signature")
			return
		}

		//if payment, _ := db.Get(transactionId).Result(); payment == "" {
		products := ProductList[productId] // Product{}
		if len(products) == 0 {
			response = errorResponse("1001", "CALLBACK_INVALID_PAYMENT", fmt.Sprintf("Wrong callback params. Product not found", productId, amount, sig, uid, transactionId))
			return
		}
		transaction := Transaction{UserId: uid, TransactionId: transactionId, ItemId: productId, Complete: false}
		err := transaction.Create(sql)
		if err != nil {
			response = errorResponse("1001", "CALLBACK_INVALID_PAYMENT", fmt.Sprintf("Wrong callback params.", err.Error()))
			return

		} else {
			response = successResponse()
		}
	}
}

func updateUserStats(message *NatsMessage, client *redis.Client) error {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(message.Message.(string)), &data)
	if len(data) == 0 {
		return nil
	}
	if message.Type == NATS_SET_MESSAGE {
		for k, v := range data {
			switch k {
			case "coins", "crystals", "candies":
				log.Printf(">>> GOT", k)
				var val float64
				val, ok := v.(float64)
				if !ok {
					intVal, _ := strconv.Atoi(v.(string))
					val = float64(intVal)
				}
				log.Println("SET ", k, val)

				client.HSet(REDIS_STATS_KEY+message.UserId, k, val).Result()
			}
		}
	}

	if message.Type == NATS_APPEND_MESSAGE {
		for k, v := range data {
			switch k {
			case "coins", "crystals", "candies":
				client.HIncrByFloat(REDIS_STATS_KEY+message.UserId, k, v.(float64)).Result()
			}
		}
	}
	return err
}

func main() {
	dbConnection := redis.NewClient(&redis.Options{
		Addr: ENV("REDIS_HOST"), //"localhost:6379",
		DB:   0,                 // use default DB
	})

	natsUrl := ENV("NATS_HOST")
	var natsConnection *nats.Conn
	var natsService *nats.EncodedConn
	var sqlDb *sql.DB
	var err error

	if natsConnection, err = nats.Connect(natsUrl); err != nil {
		log.Panicf("nats host is not available: ", err.Error())
	} else if natsService, err = nats.NewEncodedConn(natsConnection, nats.JSON_ENCODER); err != nil {
		log.Panicf("nats connecton is not available: ", err)
	} else if _, err := natsService.Subscribe(USER_INFO_CHANNEL, func(subj string, message *NatsMessage) {
		updateUserStats(message, dbConnection)
	}); err != nil {
		log.Panic("nats subscription error", err.Error())
	}

	if sqlDb, err = sql.Open("mysql", fmt.Sprintf("%s:@tcp(%s)/%s", ENV("MYSQL_USER"), ENV("MYSQL_HOST"), ENV("MYSQL_DATABASE"))); err != nil {
		log.Panic("mysql connection error", err.Error())
	}


	host := ENV("HTTP_HOST")
	http.HandleFunc("/get/",	  1receive(dbConnection, sqlDb, natsService))
	http.HandleFunc("/callback/", callback(dbConnection, sqlDb, natsService))
	http.HandleFunc("/healthcheck/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("alive"))
	})
	log.Printf("Bank Service is ready")
	log.Fatalf(http.ListenAndServe(host, nil).Error())
	for {

	}
	// http.HandleFunc("/get",   get() )
}

*/

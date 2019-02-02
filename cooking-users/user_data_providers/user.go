package user_data_providers

type User struct {
	Id       int
	Token    string
	UserId   string
	Username string
	Params   map[string]interface{}
}

const (
	DEFAULT_CANDIES_COUNT = "5"
	DEFAULT_CRYSTAL_COUNT = "100"
	DEFAULT_COINS_COUNT   = "3000"
)

func (u *User) Fill(data map[string]string) {
	u.Params = make(map[string]interface{})
	for key, value := range data {
		u.Params[key] = value
	}
}

func (u *User) SetDefault() {
	//u.Params = make(map[string]interface{})
	//u.Params["token"] = u.Token
	//u.Params["user_rating"] = "0"
	//u.Params["user_id"] = u.UserId
	//u.Params["username"] = u.Username
	//u.Params["coins"] = DEFAULT_COINS_COUNT
	//u.Params["crystals"] = DEFAULT_CRYSTAL_COUNT
	//u.Params["candies"] = DEFAULT_CANDIES_COUNT

	u.Params = map[string]interface{}{


		"token":                    u.Token,
		"user_rating":              0,
		"user_id":                  u.UserId,
		"username":                 u.Username,
		//"coins":                    DEFAULT_COINS_COUNT,
		//"crystals":                 DEFAULT_CRYSTAL_COUNT,
		"candies":                  DEFAULT_CANDIES_COUNT,

		"bun_baked_level": 1,
		"baker_level": 1,
		"bun_platform_level": 1,
		"topping_chocolate_level": 1,
		"topping_white_level": 0,
		"topping_rainbow_level": 0,
		"topping_strawberry_level": 0,
		"coffee_level": 1,
		"coffee_machine_level": 1,
		"candy": DEFAULT_CANDIES_COUNT,
		"stars": 0,
		"coins": DEFAULT_COINS_COUNT,
		"crystals": DEFAULT_CRYSTAL_COUNT,
		"bakery_level": 0,

	}
}

func (u *User) Update(params map[string]interface {
}) {
	for key, value := range params {
		u.Params[key] = value
	}
}

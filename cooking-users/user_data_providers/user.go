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
		"coins":                    DEFAULT_COINS_COUNT,
		"crystals":                 DEFAULT_CRYSTAL_COUNT,
		"candies":                  DEFAULT_CANDIES_COUNT,
		"bun_level":                1,
		"bake_level":               1,
		"bun_keeper_level":         1,
		"bakery_lamps_level":       1,
		"bakery_paintings_level":   1,
		"bakery_tables_level":      1,
		"topping_chocolate_level":  1,
		"topping_white_level":      1,
		"topping_rainbow_level":    1,
		"topping_strawberry_level": 1,
		"coffee_level":             1,
		"coffee_machine_level":     1,
		"bakery_level":             0,
		"bakery_level_1":           0,
		"bakery_level_2":           0,
		"bakery_level_3":           0,
		"bakery_level_4":           0,
		"bakery_level_5":           0,
		"bakery_level_6":           0,
		"bakery_level_7":           0,
		"bakery_level_8":           0,
		"bakery_level_9":           0,
		"bakery_level_10":          0,
		"bakery_level_11":          0,
		"bakery_level_12":          0,
		"bakery_level_13":          0,
		"bakery_level_14":          0,
		"bakery_level_15":          0,
		"bakery_level_16":          0,
		"bakery_level_17":          0,
		"bakery_level_18":          0,
	}
}

func (u *User) Update(params map[string]interface {
}) {
	for key, value := range params {
		u.Params[key] = value
	}
}

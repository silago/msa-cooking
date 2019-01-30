package user_data_providers

type UserDataProvider interface {
	DeleteAll  ( userId string )   error
	IncrementOne  (userId string, resourceName string, value interface{}) ( string, error)
	GetOne (userId string , name string ) string
	IncrementMany (userId string, values map[string]interface{} ) error
	SetMany (userId string,    values map[string]interface{} ) error
	SetOne (userId string,    key string, value interface{}) error
	GetAll (userId string /*, values map[string]float64 */ ) map[string]string
	GetOrCreate(userId string , values map[string]interface{}) (User, error)
}


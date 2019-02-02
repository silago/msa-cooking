package user_data_providers

type UserDataProvider interface {
	DeleteAll  ( userId string )   error
	IncrementOne  (userId string, resourceName string, value interface{}) ( string, error)
	IncrementMany (userId string, values map[string]interface{} ) error
	DecrementOne  (userId string, resourceName string, value interface{}) ( string, error)
	//DecrementMany (userId string, values map[string]interface{} ) error

	GetOne (userId string , name string ) string
	SetMany (userId string,    values map[string]interface{} ) error
	SetOne (userId string,    key string, value interface{}) error
	GetAll (userId string /*, values map[string]float64 */ ) map[string]string
	GetMany (userId string , names []string) map[string]string
	GetOrCreate(userId string , values map[string]interface{}) (User, error)
}


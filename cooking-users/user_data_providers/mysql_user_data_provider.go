package user_data_providers

import (
	//. "cooking-users/user_service"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

//type UserDataProvider interface {
//	getOne  (userId string, resourceName string) ( float64 , error)
//	getMany (userId string ) ( map[string]float64 , error)
//	setOne  (userId string, resourceName string, value float64)   error
//	setMany (userId string, values map[string]float64 ) error
//}

type MysqlProvider struct {
	connection *sql.DB
	tablename  string
}

func (provider MysqlProvider) SetOne(userId string, key string, value interface{}) error {
	_,err:=provider.connection.Exec("update %s set %s = '%s' ", provider.tablename, key, value);
	return  err
}

func (provider MysqlProvider) GetOne(userId string, resourceName string) string  {
	query:= fmt.Sprintf("select from {%s} where user_id = {%s} and resource_type = {%s} limit 1" , provider.tablename, userId, resourceName);
	var val string;
	if results, err := provider.connection.Query ( query); err!=nil {
	} else {
		for results.Next() {
			if err:= results.Scan(val); err!=nil {
				log.Printf(err.Error())
			}
		}
	}
	return val
}

func (provider MysqlProvider) getOneFloat64 (userId string, resourceName string) ( float64 , error) {
	query:= fmt.Sprintf("select from {%s} where user_id = {%s} and resource_type = {%s} limit 1" , provider.tablename, userId, resourceName);
	if results, err := provider.connection.Query ( query); err!=nil {
		return 0, err
	} else {
		for results.Next() {
			var val float64;
			err = results.Scan(val);
			return val, err;
		}
		return 0, nil
	}
}

func (provider MysqlProvider) GetAll (userId string )  map[string]string  {
	query:= fmt.Sprintf("select recource_name, value from {%s} where user_id = {%s} " , provider.tablename, userId);
	if results, err := provider.connection.Query ( query); err!=nil {
		return nil
	} else {
		result:= make(map[string]string)
		for results.Next() {
			var name string;
			var val  string;
			err = results.Scan(name,  val)
			result[name] = val
		}
		return result
	}
}


func (provider MysqlProvider) IncrementOne  (userId string, resourceName string, value interface{}) (string, error) {
	query:= fmt.Sprintf("update {%s} set value = {%s} values (where user_id = {%s} " , provider.tablename,
		value, userId)
	if _, err := provider.connection.Query ( query); err!=nil {
		return "0", err
	} else {
		return provider.GetOne(userId, resourceName), nil
	}
}

func (provider MysqlProvider) setOne  (userId string, resourceName string, value float64) (error) {
	query:= fmt.Sprintf("insert into {%s} (user_id, resource_name, value) values (where user_id = {%s} " , provider.tablename, userId);
	if results, err := provider.connection.Query ( query); err!=nil {
		return err
	} else {
		result:= make(map[string]float64)
		for results.Next() {
			var name string;
			var val  float64;
			err = results.Scan(name,  val)
			result[name] = val
		}
		return  nil
	}
}

func (provider MysqlProvider) SetMany (userId string, values map[string]interface{} )  error {
	if ctx, e:=provider.connection.Begin(); e!=nil {
		return e
	} else {
		for k,v := range values {
			if _, err:= ctx.Exec("update %s set %s = '%s' ", provider.tablename, k, v); err!=nil {
				if rb_error:= ctx.Rollback(); rb_error!=nil {
					return rb_error
				}
				return err
			}
		}
		return ctx.Commit()
	}
}
func (provider MysqlProvider) IncrementMany (userId string, values map[string]interface{} )  error {
	if ctx, e:=provider.connection.Begin(); e!=nil {
		return e
	} else {
		for k,v := range values {
			if _, err:= ctx.Exec("update %s set %s = %s+'%s' ", provider.tablename, k, k, v); err!=nil {
				if rb_error:= ctx.Rollback(); rb_error!=nil {
					return rb_error
				}
				return err
			}
		}
		return ctx.Commit()
	}
}


func (provider MysqlProvider) DeleteAll  ( userId string ) error {
	query:= fmt.Sprintf("delete from {%s} where user_id = {%s} " , provider.tablename, userId);
	if _, err := provider.connection.Query ( query); err!=nil {
		return err
	}
	return nil
}

func (provider MysqlProvider) GetOrCreate(userId string, params map[string]interface{}) (User, error) {
	user := User {
		UserId: userId,
	}

	if values := provider.GetAll(userId); len(values)!=0 {
		user.Fill(values)
	} else {
		user.SetDefault();
		if e:=provider.SetMany(userId, user.Params); e!=nil { return user, e }
	}

	if len(params)!=0 {
		user.Update(params)
		return user, provider.SetMany(userId,params)
	}
	return user, nil
}

func NewMysqlProvider(connectionString string, dbName string) UserDataProvider {
	//127.0.0.1:3306
	//test
	if db, err := sql.Open("mysql", "username:password@tcp(" + connectionString + ")/"+dbName); err!=nil {
		log.Println(err.Error())
		return nil
	}  else {
		provider:= MysqlProvider{db, "users_resources"}
		//provider:= &MysqlProvider{db}
		return provider
	}
}





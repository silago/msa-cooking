package user_data_providers

import (
	//. "cooking-users/user_service"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"strings"
)


type MysqlProvider struct {
	connection *sql.DB
	tablename  string
}

type tableData string
type TD = tableData
const (
	USER     tableData = "user_id"
	RESOURCE tableData = "resource_name"
	VALUE    tableData = "value"
)

//query:= fmt.Sprintf("insert into {%s} (user_id, resource_name, value) values (where user_id = {%s} " , provider.tablename, userId);

func (provider MysqlProvider) GetMany(userId string, names []string) map[string]string {
	resources:=strings.Join(names,",")
	query:= fmt.Sprintf("select %s, %s  from {%s} where  = {%s} and resource_type in ({%s}) limit 1" , RESOURCE, VALUE, provider.tablename, userId, resources);
	result:=make(map[string]string)
	resource:=struct {
		name string
		value string
	}{}

	if results, err := provider.connection.Query ( query); err!=nil {
	} else {
		for results.Next() {
			if err:= results.Scan(&resource.name, &resource.value); err!=nil {
				result[resource.name]=resource.value
			}
		}
	}
	return result
}

func (provider MysqlProvider) DecrementOne(userId string, resourceName string, value interface{}) (string, error) {
	var val int
	switch value.(type) {
		case int:
			val = value.(int)
			break
		case string:
			val, _ = strconv.Atoi(value.(string))
	}
	if currentValue, e:=strconv.Atoi(provider.GetOne(userId, resourceName)); e!=nil {
		return "",e
	} else {
		currentValue:=currentValue-val
		e = provider.SetOne(userId, resourceName , string(currentValue - val))
		return string(currentValue), e
	}
}


func (provider MysqlProvider) CreateOne(userId string, key string, value interface{}) error {
	queryString:= fmt.Sprintf("insert into %s (%s, %s, %s) values () ", provider.tablename, VALUE, USER, RESOURCE)
	_,err:=provider.connection.Exec	(queryString, value,  userId, key)
	return  err
}

func (provider MysqlProvider) SetOne(userId string, key string, value interface{}) error {
	queryString:= fmt.Sprintf("update %s set %s = ? where %s = ? and %s = ? limit 1", provider.tablename, VALUE, USER, RESOURCE)
	if result,err:=provider.connection.Exec	(queryString, value,  userId, key); err!=nil {
		return err
	} else if count,_:=result.RowsAffected(); count==0 {
		return provider.CreateOne(userId, key, value)
	}
	return nil
}

func (provider MysqlProvider) GetOne(userId string, resourceName string) string  {
	query:= fmt.Sprintf("select from %s where %s = ? and %s = ? limit 1" , provider.tablename, USER, RESOURCE)
	var val string
	if results, err := provider.connection.Query ( query, userId, resourceName); err!=nil {
	} else {
		for results.Next() {
			if err:= results.Scan(val); err!=nil {
				log.Printf(err.Error())
			}
		}
	}
	return val
}

func (provider MysqlProvider) GetAll (userId string )  map[string]string  {
	query:= fmt.Sprintf("select %s, %s from %s where user_id = ? " , RESOURCE, VALUE, provider.tablename);
	if results, err := provider.connection.Query ( query, userId); err!=nil {
		return nil
	} else {
		result:= make(map[string]string)
		for results.Next() {
			var name string
			var val  string
			err = results.Scan(name,  val)
			result[name] = val
		}
		return result
	}
}


func (provider MysqlProvider) IncrementOne  (userId string, resourceName string, value interface{}) (string, error) {
	query:= fmt.Sprintf("update %s set %s = ? where  %s = ? and %s =? limit 1 " , provider.tablename,VALUE, USER, RESOURCE)
	if _, err := provider.connection.Query ( query, value, userId, resourceName); err!=nil {
		return "0", err
	} else {
		return provider.GetOne(userId, resourceName), nil
	}
}


func (provider MysqlProvider) SetMany (userId string, values map[string]interface{} )  error {
	if ctx, e:=provider.connection.Begin(); e!=nil {
		return e
	} else {
		for k,v := range values {
			if _, err:= ctx.Exec(fmt.Sprintf("update %s set %s = ? where %s = ? and %s = ? ", provider.tablename, VALUE, USER, RESOURCE),v,userId,k); err!=nil {
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
			if result, err:= ctx.Exec(fmt.Sprintf("update %s set %s = %s + ? where %s = ? and %s = ? ", provider.tablename, VALUE, VALUE, USER, RESOURCE),v,userId,k); err!=nil {
				if rb_error:= ctx.Rollback(); rb_error!=nil {
					return rb_error
				}
				return err
			} else if updated, _:=result.RowsAffected(); updated ==0 {
				if e:=provider.SetOne(userId, k, v); e!=nil {
					log.Println(e.Error())
				}
			}
		}
		return ctx.Commit()
	}
}


func (provider MysqlProvider) DeleteAll  ( userId string ) error {
	query:= fmt.Sprintf("delete from %s where user_id = ? " , provider.tablename)//, userId);
	if _, err := provider.connection.Query ( query, userId); err!=nil {
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
	if db, err := sql.Open("mysql", "username:password@tcp(" + connectionString + ")/"+dbName); err!=nil {
		log.Println(err.Error())
		return nil
	}  else {
		provider:= MysqlProvider{db, "users_resources"}
		return provider
	}
}





//插入数据 or 查询数据
package main

import ( 
	"database/sql"
	_ "github.com/xxiaomo/mssql"
	"net/http"
	"fmt"
	"gopkg.in/couchbase/gocb.v1"
	"crypto/md5"	
	"encoding/json"
)

//用户信息
type User struct {
	Id		int 	`json:"id"`
	Name	string	`json:"name"` 
	Age		int 	`json:"age"`
}

//Insert into database
func insertIntoDB(w http.ResponseWriter, r *http.Request) {
	//构造查询字符串
	r.ParseForm()	//解析参数
	nameValue := r.FormValue("name")	//用户姓名
	ageValue := r.FormValue("age")		//用户年龄
	// if nameValue == "" || ageValue == ""{	//表单检查
	// 	fmt.Fprintf(w, "Usage:/insert?name=xx&age=xx")
	// 	return 
	// }
	query := fmt.Sprintf("INSERT INTO users(name, age) VALUES('%s', '%s')", nameValue, ageValue)	

	//打开数据库连接池
	//格式 sqlserver://sa:123456@localhost:1433?database=users&connection+timeout=30"
	//	sa=username, 123456=password, 1433=port
	db, _ := sql.Open("sqlserver", "sqlserver://sa:123456@localhost?database=users&connection+timeout=30")

	//往数据库中插入数据
	_, err := db.Exec(query)	//执行INSERT语句
	if err != nil {
		fmt.Fprintf(w, "Insert fail")
	} else {
		fmt.Fprintf(w, "Insert success")
	}
}

//Select from database
func selectFromDB(w http.ResponseWriter, r *http.Request) {
	//构造查询语句
	r.ParseForm()
	if(len(r.Form) != 1) {
		fmt.Fprintf(w, "Usage: /select?key=value")
		return 
	}
	var query string
	for k, _ := range r.Form {
		query = fmt.Sprintf("SELECT id, name, age FROM users WHERE %s = '%s'", k, r.FormValue(k))
	}

	//生成key值
	key := fmt.Sprintf("%x", md5.Sum([]byte(query)))

	//连接到couchbase缓存桶
	cluster, _ := gocb.Connect("couchbase://localhost")
	bucket, _ := cluster.OpenBucket("userdata", "")

	//到缓存桶中查找数据
	var data string
	_, err := bucket.Get(key, &data)
	if err == nil { //在缓存中找到数据
		//输出数据并退出
		fmt.Fprintf(w, "%s", data)
		return 
	}

	//缓存中没有数据,则到sqlserver中查找数据

	//打开数据库连接池
	//格式 sqlserver://sa:123456@localhost:1433?database=users&connection+timeout=30"
	//	sa=username, 123456=password, 1433=port
	db, _ := sql.Open("sqlserver", "sqlserver://sa:123456@localhost?database=users&connection+timeout=30")

	rows, err := db.Query(query)	//执行Query语句
	//声明数据结构存储查询数据
	var userdata []User
	var temp User
	for rows.Next() {	//遍历查询结果
		rows.Scan(&temp.Id, &temp.Name, &temp.Age)
		userdata = append(userdata, temp)
	}
	if len(userdata) == 0 {	//没有在数据库中查到数据
		fmt.Fprintf(w, "not found")
		return 
	}
	//生成json字符串
	jsondata , _ := json.Marshal(userdata)
	data = string(jsondata)
	//把数据插入缓存
	bucket.Upsert(key , data, 0)

	//输出数据
	fmt.Fprintf(w, "%s", data)
}

func main() {
	http.HandleFunc("/insert", insertIntoDB)
	http.HandleFunc("/select", selectFromDB)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err.Error())
		return
	}
}
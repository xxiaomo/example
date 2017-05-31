//插入数据 or 查询数据
package main

import ( 
	"database/sql"
	_ "github.com/xxiaomo/mssql"
	"net/http"
	"fmt"
)

//Insert into sql server
func insertIntoServer(w http.ResponseWriter, r *http.Request) {
	//构造查询字符串
	r.ParseForm()	//解析参数
	nameValue := r.FormValue("name")	//得到键值为“name”的参数的值
	query := fmt.Sprintf("INSERT INTO users(name) VALUES('%s')", nameValue)	

	//打开数据库连接池
	//格式 sqlserver://sa:123456@localhost:1433?database=users&connection+timeout=30"
	//	sa=username, 123456=password, 1433=port
	db, err := sql.Open("sqlserver", "sqlserver://sa:123456@localhost?database=users&connection+timeout=30")
	if err != nil {
		fmt.Println("sql open fail")
		return
	}
	defer db.Close()

	//往数据库中插入数据
	_, err = db.Exec(query)	//执行INSERT语句
	if err != nil {
		fmt.Fprintf(w, "Insert fail")
	} else {
		fmt.Fprintf(w, "Insert success")
	}
}

//Select from sql server
func selectFromServer(w http.ResponseWriter, r *http.Request) {
	//构造查询字符串
	r.ParseForm()	//解析参数
	nameValue := r.FormValue("name")	//得到键值为“name”的参数的值
	query := fmt.Sprintf("SELECT id, name FROM users WHERE name = '%s'", nameValue)

	//打开数据库连接池
	//格式 sqlserver://sa:123456@localhost:1433?database=users&connection+timeout=30"
	//	sa=username, 123456=password, 1433=port
	db, err := sql.Open("sqlserver", "sqlserver://sa:123456@localhost?database=users&connection+timeout=30")
	if err != nil {
		fmt.Println("sql open fail")
		return
	}
	defer db.Close()

	//查询数据
	var (
		id int
		name string
	)
	rows, err := db.Query(query)	//执行Query语句
	if err != nil {		
		fmt.Fprintf(w, "select fail")
	} else {	//遍历并输出查询结果
		for rows.Next() { 	
			err = rows.Scan(&id, &name)
			fmt.Fprintln(w, id, name)
		}	
	}
}

func main() {
	http.HandleFunc("/insert", insertIntoServer)
	http.HandleFunc("/select", selectFromServer)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err.Error())
		return
	}
}
















// //sql server 插入查询示例，为了简洁，示例不会进行错误检查
// package main

// import (
// 	"database/sql"
// 	_ "github.com/xxiaomo/mssql"
// 	"fmt"	//Package fmt implements formatted I/O
// )

// func main() {
// 	//打开数据库连接池
// 	//格式 sqlserver://sa:123456@localhost:1433?database=users&connection+timeout=30"
// 	//	sa=username, 123456=password, 1433=port
	// db, err := sql.Open("sqlserver", "sqlserver://sa:123456@localhost?database=users&connection+timeout=30")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer db.Close()

// 	//往数据库中插入数据
// 	_, err = db.Exec("INSERT INTO users(name) VALUES('xdy')")	//执行INSERT语句
// 	if err != nil {
// 		fmt.Println("Insert fail")
// 	} else {
// 		fmt.Println("Insert success")
// 	}

// 	//查询刚才插入的数据
// 	var (
// 		id int
// 		name string
// 	)
// 	rows, err := db.Query("SELECT id, name FROM users WHERE id = 1")	//执行Query语句
// 	for rows.Next() { 	//遍历并输出查询结果
// 		err = rows.Scan(&id, &name)
// 		fmt.Println(id, name)
// 	}
// }
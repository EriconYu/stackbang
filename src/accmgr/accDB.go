package accmgr

import (
	"database/sql"
	"fmt"
	"time"
	"randomlib"

	_ "github.com/Go-SQL-Driver/MySQL"
)

// UserBasicInfoDB 用户基础信息
type UserBasicInfoDB struct {
	UserID      uint64 `json:"UserID"`
	UserName    string `json:"UserName"`
	Password    string `json:"Password"`
	Email       string `json:"Email"`
	PhoneNumber string `json:"PhoneNumber"`
}

// LoadByEmail 登录
func (u *UserBasicInfoDB) LoadByEmail() (AccessToken string ,ok bool) {
	// sqlinit 这个库必须重新打开才能使用
	db, e := sql.Open("mysql", "root:4759694294@tcp(stackbang.com:3306)/userInfo?charset=utf8")
	defer db.Close()
	if e != nil {
		fmt.Println("mysql error")
	}

	cmd := fmt.Sprintf(`SELECT UserID , UserName,Password,Email,PhoneNumber FROM userLogin where Email = "%s";`, u.Email)

	rows, err := db.Query(cmd)
	if err != nil {
		fmt.Println("LoadByEmail err is ", err)
		return "",false
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&u.UserID, &u.UserName, &u.Password, &u.Email, &u.PhoneNumber)
		if err != nil {
			fmt.Println("rows Scan err ! error code is :", err)
			return "",false
		}
	}
	AccessToken = string(randomlib.Krand(32,randomlib.KC_RAND_KIND_ALL))
	return AccessToken , true

}

// InsertUser 注册
func (u *UserBasicInfoDB) InsertUser() (ok bool) {
	db, e := sql.Open("mysql", "root:4759694294@tcp(stackbang.com:3306)/userInfo?charset=utf8")
	defer db.Close()
	if e != nil {
		fmt.Println("mysql error")
	}
	//插入一条新纪录
	stmt, err := db.Prepare("insert into userLogin(UserName,Password,Email,PhoneNumber,RegistTime)values(?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	RegistTime := fmt.Sprintf("%s", time.Now())
	result, err := stmt.Exec(u.UserName, u.Password, u.Email, u.PhoneNumber, RegistTime)
	if err == nil {
		id, err := result.LastInsertId()
		if err == nil {
			u.UserID = uint64(id)
			return true
		}
	}
	fmt.Println("err is ", err)
	return true
}

// IsExistEmail 邮箱查重 true 存在  false 不存在
func (u *UserBasicInfoDB) IsExistEmail() bool {
	// sqlinit 这个库必须重新打开才能使用
	db, e := sql.Open("mysql", "root:4759694294@tcp(stackbang.com:3306)/userInfo?charset=utf8")
	defer db.Close()
	if e != nil {
		fmt.Println("mysql error")
	}
	// 查询数据
	rows, err := db.Query("select count(Email) from userLogin where Email=\"" + u.Email + "\";")

	if err != nil {
		fmt.Println("exec failurr! err is ", err)
		return false
	}
	if rows == nil {
		fmt.Println("rows is null")
		return false
	}
	var count int
	rows.Next()
	rows.Scan(&count)
	if count == 0 {
		return false
	}
	return true
}

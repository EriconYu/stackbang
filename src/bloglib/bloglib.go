package bloglib

import (
	"database/sql"
	"fmt"

	_ "github.com/Go-SQL-Driver/MySQL"
)

var db *sql.DB

// Blog Blog结构体
type Blog struct {
	ID         uint64 `json:"ID"`
	AuthorID   uint64 `json:"AuthorID"`
	AuthorName string `json:"AuthorName"`
	Date       string `json:"Date"`
	Title      string `json:"Title"`
	Content    string `json:"Content"`
	KeyWords   string `json:"KeyWords"`
	Classify   string `json:"Classify"`
	BlogURL    string `json:"BlogURL"`
}

// ReleaseBlog 发布一篇博客
func ReleaseBlog(b Blog) (ok bool) {
	db, _ := sql.Open("mysql", "root:4759694294@tcp(stackbang.com:3306)/userBlog?charset=utf8")
	defer db.Close()

	// 插入数据
	Str := fmt.Sprintf(`insert into Blogs(AuthorID ,AuthorName ,Date , Title , Content ,KeyWords ,Classify ,BlogURL) VALUES (%d , %s , %s , %s , %s , %s , %s,%s );`, b.AuthorID, b.AuthorName, b.Date, b.Title, b.Content, b.KeyWords, b.Classify, b.BlogURL)
	_, e5 := db.Query(Str)
	if e5 == nil {
		fmt.Println("insert table success!")
		return true
	}
	fmt.Println("insert table failer ", e5)
	return false
}

// GetBlogs ..
func GetBlogs() (BlogCom []Blog) {
	var e error
	var BlogStTemp Blog
	// sqlinit 这个库必须重新打开才能使用
	db, e = sql.Open("mysql", "root:4759694294@tcp(stackbang.com:3306)/userBlog?charset=utf8")
	defer db.Close()
	if e != nil {
		fmt.Println("mysql error , e is ", e)
	}

	rows, _ := db.Query("SELECT ID, AuthorID ,AuthorName ,Date , Title , Content ,KeyWords ,Classify ,BlogURL FROM Blogs ORDER BY Date ;")
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&BlogStTemp.ID, &BlogStTemp.AuthorID, &BlogStTemp.AuthorName, &BlogStTemp.Date, &BlogStTemp.Title, &BlogStTemp.Content, &BlogStTemp.KeyWords, &BlogStTemp.Classify, &BlogStTemp.BlogURL)
		if err != nil {
			fmt.Println("rows Scan err ! error code is :", err)
		}
		BlogCom = append(BlogCom, BlogStTemp)
	}
	return BlogCom
}

// GetBlog ..
func GetBlog(BlogID uint64) (BlogCom Blog, ok bool) {
	var e error
	// sqlinit 这个库必须重新打开才能使用
	db, e = sql.Open("mysql", "root:4759694294@tcp(stackbang.com:3306)/userBlog?charset=utf8")
	defer db.Close()
	if e != nil {
		fmt.Println("mysql error , e is ", e)
	}

	cmd := fmt.Sprintf("SELECT * FROM Blogs where ID = %d ;", BlogID)

	rows, err := db.Query(cmd)
	if err != nil {
		fmt.Println("GetBlog err is ", err)
		return BlogCom, false
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&BlogCom.ID, &BlogCom.AuthorID, &BlogCom.AuthorName, &BlogCom.Date, &BlogCom.Title, &BlogCom.Content, &BlogCom.KeyWords, &BlogCom.Classify, &BlogCom.BlogURL)
		if err != nil {
			fmt.Println("rows Scan err ! error code is :", err)
			return BlogCom, false

		}
	}
	fmt.Println("BlogCom is ", BlogCom)
	return BlogCom, true
}

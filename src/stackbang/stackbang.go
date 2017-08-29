package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"accmgr"
	"bloglib"
)

// loginFunc 登录
func loginFunc(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("loginFunc method:", req.Method) //获取请求的方法
	if req.Method == "GET" {
		t, _ := template.ParseFiles("../www/login.html")
		t.Execute(rw, nil)
	} else if req.Method == "POST" {
		//请求的是登录数据，那么执行登录的逻辑判断
		req.ParseForm()
		fmt.Println("inputEmail:", req.Form["inputEmail"])
		fmt.Println("inputPassword:", req.Form["inputPassword"])
		var lenEmail = len(req.Form["inputEmail"])
		if lenEmail == 0 {
			http.Redirect(rw, req, "/login", 302)
			return
		}
		var lenPwd = len(req.Form["inputPassword"])
		if lenPwd == 0 {
			http.Redirect(rw, req, "/login", 302)
			return
		}
		var email = req.FormValue("inputEmail")
		var pwd = req.FormValue("inputPassword")
		if 0 == len(email) {
			http.Redirect(rw, req, "/login", 302)
			return
		}
		var userAccMgr accmgr.UserAccMgr
		userAccMgr.Email = email
		userAccMgr.Password = pwd
		var ok = userAccMgr.LoginByEmail()
		if ok == true { //登录成功 写入cookie
			fmt.Println("login success")
			//写入cookie
			expiration := time.Now()
			expiration = expiration.AddDate(1, 0, 0)
			cookie := http.Cookie{Name: "stackbangusername", Value: userAccMgr.UserName, Expires: expiration, Path: "/"}
			http.SetCookie(rw, &cookie)
			cookie = http.Cookie{Name:"stackbanguserid",Value:fmt.Sprintf("%d",userAccMgr.UserID) , Expires: expiration,Path:"/"}
			http.SetCookie(rw, &cookie)
			cookie = http.Cookie{Name:"stackbangusertoken",Value:userAccMgr.AccessToken , Expires: expiration,Path:"/"}
			http.SetCookie(rw, &cookie)
			http.Redirect(rw, req, "/index", 302)
		} else { //登录失败，仍然在登录页面
			fmt.Println("login failer")
			http.Redirect(rw, req, "/login", 302)
		}
	}
}

// logoutFunc ..
func logoutFunc(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("logoutFunc method:", req.Method) //获取请求的方法

}

// registFunc 注册账户
func registFunc(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("registFunc method:", req.Method) //获取请求的方法
	if req.Method == "GET" {
		t, _ := template.ParseFiles("../www/regist.html")
		t.Execute(rw, nil)
	} else if req.Method == "POST" {
		req.ParseForm()
		var regsitinfo accmgr.UserAccMgr
		regsitinfo.UserName = req.FormValue("inputName")
		regsitinfo.Email = req.FormValue("inputEmail")
		regsitinfo.Password = req.FormValue("inputPassword")
		ConfirmPassword := req.FormValue("inputConfirmPassword")
		if regsitinfo.Password != ConfirmPassword {
			http.Redirect(rw, req, "/regist/", 302)
			return
		}
		if regsitinfo.UserName == "" || regsitinfo.Email == "" || regsitinfo.Password == "" {
			http.Redirect(rw, req, "/regist/", 302)
			return
		}
		registOk := regsitinfo.Regist()
		fmt.Println("The registOk is ", registOk)
		if registOk == 0 {
			http.Redirect(rw, req, "/login", 302)
		} else {
			http.Redirect(rw, req, "/regist", 302)
		}
	}

}

// indexFunc 首页
func indexFunc(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("indexFunc method:", req.Method) //获取请求的方法
	fmt.Println("indexFunc Remote is", req.Header.Get("X-Real-IP"))
	if req.Method == "GET" {
		t, _ := template.ParseFiles("../www/index.html")
		t.Execute(rw, nil)
	}
}

// writeBlogFunc ..
func writeBlogFunc(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("writeBlogFunc method:", req.Method) //获取请求的方法
	if req.Method == "GET" {
		t, _ := template.ParseFiles("../www/blog/writeblog.html")
		t.Execute(rw, nil)
	} else if req.Method == "POST" {
		var blog bloglib.Blog
		req.ParseForm()
		//		cookieAuthor , err := req.Cookie("StackBangUserName")
		//		if err == nil{
		//			fmt.Println("writeBlogFunc POST blog cookie err is :\n" , err)
		//			fmt.Fprintln(rw, "cookie error!")
		//		}
		//		blog.AuthorName = cookieAuthor.String()
		blog.AuthorName = "eric"
		blog.Content = req.FormValue("content")
		blog.Date = fmt.Sprintf("%d-%d-%d %d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
		if true == bloglib.ReleaseBlog(blog) {
			http.Redirect(rw, req, "/blogindex", 200)
		} else {
			http.Redirect(rw, req, "/", 200)
		}
	}

}

// blogIndexFunc blog首页
func blogIndexFunc(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("blogIndexFunc method:", req.Method) //获取请求的方法
	if req.Method == "GET" {
		t, _ := template.ParseFiles("../www/blog/blogindex.html")
		t.Execute(rw, nil)
	}
}

// getBlogsFunc 获取博客内容
func getBlogsFunc(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("getBlogsFunc method:", req.Method) //获取请求的方法
	fmt.Println("getBlogsFunc Remote is", req.Header.Get("X-Real-IP"))
	var Blogs []bloglib.Blog
	if req.Method == "GET" {
		Blogs = bloglib.GetBlogs()
		BlogsJSON, err := json.Marshal(Blogs)
		if err != nil {
			fmt.Println("Json Marshal err is ", err)
			return
		}
		fmt.Fprintf(rw, string(BlogsJSON))
	}
}

// aboutmeFunc 关于我
func aboutmeFunc(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("aboutmeFunc method:", req.Method) //获取请求的方法
	fmt.Println("aboutmeFunc Remote is", req.Header.Get("X-Real-IP"))
	if req.Method == "GET" {
		t, _ := template.ParseFiles("../www/blog/aboutme.html")
		t.Execute(rw, nil)
	}
}

// getBlogFunc 获取博客内容
func getBlogFunc(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("blog method:", req.Method) //获取请求的方法

	if req.Method == "GET" {
		fmt.Println("req.RequestURI is ", req.RequestURI)
		BlogIDStr := req.RequestURI[12:]
		BlogID, _ := strconv.Atoi(BlogIDStr)
		Blog, ok := bloglib.GetBlog(uint64(BlogID))
		if ok != true {
			fmt.Println("bloglib.GetBlog false")
			fmt.Fprintf(rw, string("ServerError"))
		}
		BlogJSON, err := json.Marshal(Blog)
		if err != nil {
			fmt.Println("getBlogFunc Json Marshal err is ", err)
		}
		fmt.Fprintf(rw, string(BlogJSON))
	}
}

// imgFunc ..
func imgFunc(rw http.ResponseWriter, req *http.Request) {

}

// HD 过滤文件服务器访问
type HD struct {
	http.Handler
}

// ServeHTTP ..
func (h HD) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	lenPath := len(req.URL.Path)
	if req.URL.Path[lenPath-1] == '/' {
		fmt.Println("Path is illegal!!!")
		http.NotFound(rw, req)
		return
	}
	hd := http.StripPrefix("/img/", http.FileServer(http.Dir("../img/")))
	hd.ServeHTTP(rw, req)
}

func main() {
	fmt.Println("Hello StackBang!")
	var ff HD
	http.HandleFunc("/", indexFunc)               //设置访问的路由
	http.HandleFunc("/index/", indexFunc)         //设置访问的路由
	http.HandleFunc("/logout/", logoutFunc)       //设置访问的路由
//	http.HandleFunc("/login/", loginFunc)         //设置访问的路由
//	http.HandleFunc("/regist/", registFunc)       //设置访问的路由
	http.HandleFunc("/writeblog/", writeBlogFunc) //设置访问的路由
	http.HandleFunc("/blogindex/", blogIndexFunc) //设置访问的路由
	http.HandleFunc("/aboutme/", aboutmeFunc)     //设置访问的路由
	http.HandleFunc("/getblogs/", getBlogsFunc)   //设置访问的路由
	http.HandleFunc("/getblog/", getBlogFunc)     //设置访问的路由

	http.Handle("/common/", http.StripPrefix("/common/", http.FileServer(http.Dir("../common"))))

	//	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("../img"))))
	http.Handle("/img/", ff)
	err := http.ListenAndServe(":8888", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

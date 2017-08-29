package accmgr

import _ "github.com/Go-SQL-Driver/MySQL"

// UserProfile 用户详细信息
type UserProfile struct {
	UserID    uint64 `json:"UserID"`
	School    string `json:"School"`
	Address   string `json:"Address"`
	Company   string `json:"Company"`
	Sex       string `json:"Sex"`
	Job       string `json:"Job"`
	UserIntro string `json:"UserIntro"`
}

// UserAccMgr 用户账户管理
type UserAccMgr struct {
	UserID      uint64 `json:"UserID"`
	UserName    string `json:"UserName"`
	Password    string `json:"Password"`
	Email       string `json:"Email"`
	PhoneNumber string `json:"PhoneNumber"`
	AccessToken string `json:"AccessToken"`
}

// LoginByEmail 登录
func (u *UserAccMgr) LoginByEmail() (ok bool) {
	var userBasicInfoDB UserBasicInfoDB
	userBasicInfoDB.Email = u.Email
	u.AccessToken, ok = userBasicInfoDB.LoadByEmail()
	if ok == false {
		return false
	}
	if u.Password == userBasicInfoDB.Password {
		return true
	}
	return false
}

// Regist 注册
func (u *UserAccMgr) Regist() (ErrCode int) {
	var userBasicInfoDB UserBasicInfoDB
	userBasicInfoDB.Email = u.Email
	userBasicInfoDB.Password = u.Password
	userBasicInfoDB.UserName = u.UserName

	ok := userBasicInfoDB.IsExistEmail()
	if ok == true {
		return -2 // 已存在该邮箱，注册失败
	}
	ok = userBasicInfoDB.InsertUser()
	if ok == false {
		return -1 // 注册失败，数据库插入失败
	}
	return 0
}

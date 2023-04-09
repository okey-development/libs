package service

import "database/sql"

type User struct {
	Id           int64  `json:"user_id"`
	ReferralId   int64  `json:"referral_id"`
	Username     string `json:"username"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	LangId       int64  `json:"lang_id"`
	Email        string `json:"email"`
	EmailConfirm bool   `json:"email_confirm"`
	MainCard     int64  `json:"main_card"`
}

func GetUserByID(user_id int64) (*User, error) {
	row := QueryRowDB(`SELECT u.id, 
	u.referredby, 
	u.account, 
	u.firstname, 
	u.lastname, 
	u.lang, 
	u.email, 
	u.email_confirm, 
	u.main_card
	FROM admin.users u 
	where u.id= $1 and (u.is_archive = false or u.is_archive isnull);`, user_id)

	var (
		userSql UserSQL
	)
	if err := row.Scan(&userSql.Id,
		&userSql.ReferralId,
		&userSql.Username,
		&userSql.Firstname,
		&userSql.Lastname,
		&userSql.LangId,
		&userSql.Email,
		&userSql.EmailConfirm,
		&userSql.MainCard); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}

	return userSql.Scan(), nil
}

type UserSQL struct {
	Id           sql.NullInt64
	ReferralId   sql.NullInt64
	Username     sql.NullString
	Firstname    sql.NullString
	Lastname     sql.NullString
	LangId       sql.NullInt64
	Email        sql.NullString
	EmailConfirm sql.NullBool
	MainCard     sql.NullInt64
}

func (user *UserSQL) Scan() *User {
	return &User{Id: user.Id.Int64,
		ReferralId:   user.ReferralId.Int64,
		Username:     user.Username.String,
		Firstname:    user.Firstname.String,
		Lastname:     user.Lastname.String,
		LangId:       user.LangId.Int64,
		Email:        user.Email.String,
		EmailConfirm: user.EmailConfirm.Bool,
		MainCard:     user.MainCard.Int64,
	}
}

func (user *User) CheckAccesses(accesses map[Object]Action) bool {

	for object, action := range accesses {
		if !checkAccess(object, action, user.Id) {
			return false
		}
	}

	return true
}

func (user *User) CheckRights(rights map[Object]Action) bool {
	for object, action := range rights {
		if !checkRights(object, action, user.Id) {
			return false
		}
	}

	return true
}

func checkAccess(object Object, action Action, userId int64) bool {
	return true
}

func checkRights(object Object, action Action, userId int64) bool {
	return true
}

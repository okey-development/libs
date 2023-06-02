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
	Location     string `json:"location"`
	UUID         string `json:"uuid"`
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
	u.main_card,
	u.location,
	u.uuid
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
		&userSql.MainCard,
		&userSql.Location,
		&userSql.UUID,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}

	return userSql.Scan(), nil
}

func GetUserByUUID(uuid string) (*User, error) {
	row := QueryRowDB(`SELECT u.id, 
	u.referredby, 
	u.account, 
	u.firstname, 
	u.lastname, 
	u.lang, 
	u.email, 
	u.email_confirm, 
	u.main_card,
	u.location,
	u.uuid
	FROM admin.users u 
	where u.uuid= $1 and (u.is_archive = false or u.is_archive isnull);`, uuid)

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
		&userSql.MainCard,
		&userSql.Location,
		&userSql.UUID,
	); err != nil {
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
	Location     sql.NullString
	UUID         sql.NullString
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
		Location:     user.Location.String,
		UUID:         user.UUID.String,
	}
}

func (user *User) CheckAccesses(accesses map[Object][]Action) bool {

	for object, actions := range accesses {
		for _, action := range actions {
			if !checkAccess(object, action, user.Id) {
				return false
			}
		}
	}

	return true
}

func (user *User) CheckRights(rights map[Object][]Action) bool {

	for object, actions := range rights {
		for _, action := range actions {
			if !checkRights(object, action, user.Id) {
				return false
			}
		}
	}

	return true
}

func checkAccess(object Object, action Action, userId int64) bool {

	row := QueryRowDB(`select
	count(*)
	from tariff.toc_tariff_user_price ttup 
	left join tariff.toc_tariffs_accesses tta on tta.tariff_id = ttup.tariff_id 
	left join tariff.accesses a on a.id  = tta.access_id 
	left join tariff.objects o on o.id = a.obj_id 
	left join tariff.actions a2 on a2.id = a.action_id 
	where ttup.user_id = $1
	and o."name" = $2
	and a2."name" = $3
	and ttup.status = 1
	and (ttup.day_of_payment > current_date or ttup.day_of_payment isnull)
	;`, userId, object, action)

	var count sql.NullInt64
	if err := row.Scan(&count); err != nil {
		if err != sql.ErrNoRows {
			Error(Errorf(err.Error()))
			return false
		}
		return false
	}
	return count.Int64 > 0
}

func checkRights(object Object, action Action, userId int64) bool {

	row := QueryRowDB(`select
	count(*)
	from "admin".users u 
	left join roles.toc_roles_rights trr on trr.role_id = u.role_id 
	left join roles.rights r on r.id = trr.right_id  
	left join roles.objects o on o.id = r.obj_id 
	left join roles.actions a on a.id = r.action_id 
	where u.id = $1
	and o."name" = $2
	and a."name" = $3
	;`, userId, object, action)

	var count sql.NullInt64
	if err := row.Scan(&count); err != nil {
		if err != sql.ErrNoRows {
			Error(Errorf(err.Error()))
			return false
		}
		return false
	}
	return count.Int64 > 0
}

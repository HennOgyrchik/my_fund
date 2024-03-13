package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	_ "regexp"
	"time"
)

func dbConnection() (*sql.DB, error) {
	connStr := "user=postgres password=111 dbname=postgres sslmode=disable host=192.168.0.116 port=5432" //как то убрать логин и пароль, заменить ip на имя контейнера
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		_ = db.Close()
	}

	return db, err
}

func IsMember(memberId int64) (bool, error) {
	db, err := dbConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	stmt, err := db.Prepare("select count(*) from members where member_id=$1")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(memberId).Scan(&count)

	if (err != nil) || (count == 0) {
		return false, err
	}

	return true, nil
}

func IsAdmin(memberId int64) (bool, error) {
	db, err := dbConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	stmt, err := db.Prepare("select admin from members m  where member_id=$1")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var result bool
	err = stmt.QueryRow(memberId).Scan(&result)
	return result, err
}

func DoesTagExist(tag string) (bool, error) {
	db, err := dbConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	stmt, err := db.Prepare("select count(*) from funds where tag=$1")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(tag).Scan(&count)
	switch {
	case err != nil:
		return false, err
	case count > 0:
		return true, err
	default:
		return false, err
	}
}

func CreateFund(tag string, balance float64) error {
	db, err := dbConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into funds (tag,balance) values ($1,$2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_ = stmt.QueryRow(tag, balance)
	return err
}

func GetAdminFund(tag string) (int64, error) {
	db, err := dbConnection()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	stmt, err := db.Prepare("select member_id from members where tag = $1 and admin = true")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var memberId int64

	err = stmt.QueryRow(tag).Scan(&memberId)
	return memberId, err
}

// ShowBalance возвращает баланс фонда
func ShowBalance(tag string) (float64, error) {
	db, err := dbConnection()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	stmt, err := db.Prepare("select balance from funds where tag=$1")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var balance float64
	err = stmt.QueryRow(tag).Scan(&balance)
	return balance, err
}

// DeleteFund удаляет фонд
func DeleteFund(tag string) error {
	db, err := dbConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("call delete_fund($1)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_ = stmt.QueryRow(tag)

	return err
}

// GetTag возвращает тег фонда, в котором пользователь находится
func GetTag(memberId int64) (string, error) {
	db, err := dbConnection()
	if err != nil {
		return "", err
	}
	defer db.Close()

	stmt, err := db.Prepare("select tag from members where member_id=$1")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var tag string
	err = stmt.QueryRow(memberId).Scan(&tag)
	return tag, err
}

type Member struct {
	ID      int64
	Tag     string
	IsAdmin bool
	Login   string
	Name    string
}

func AddMember(member Member) error {
	db, err := dbConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into members (tag,member_id,admin,login,name) values ($1,$2,$3,$4,$5)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_ = stmt.QueryRow(member.Tag, member.ID, member.IsAdmin, member.Login, member.Name)
	return err
}

// GetMembers возвращает список пользователей фонда
func GetMembers(tag string) ([]Member, error) {
	var members []Member
	db, err := dbConnection()
	if err != nil {
		return members, err
	}
	defer db.Close()

	stmt, err := db.Prepare("select member_id, tag, admin, login, name from members where tag =$1")
	if err != nil {
		return members, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(tag)
	if err != nil {
		return members, err
	}
	defer rows.Close()

	for rows.Next() {
		var member Member
		if err = rows.Scan(&member.ID, &member.Tag, &member.IsAdmin, &member.Login, &member.Name); err != nil {
			return members, err
		}
		members = append(members, member)
	}
	return members, nil
}

// GetInfoAboutMember возвращает полную информацию о пользователе
func GetInfoAboutMember(memberId int64) (Member, error) {
	member := Member{ID: memberId}

	db, err := dbConnection()
	if err != nil {
		return member, err
	}
	defer db.Close()

	stmt, err := db.Prepare("select tag,admin,login,name from members where member_id = $1")
	if err != nil {
		return member, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(memberId).Scan(&member.Tag, &member.IsAdmin, &member.Login, &member.Name)
	return member, err
}

type CashCollection struct {
	ID         int
	Tag        string
	Sum        float64
	Status     string
	Comment    string
	CreateDate time.Time
	CloseDate  time.Time
	Purpose    string
}

func CreateCashCollection(cashCollection CashCollection) (int, error) {

	db, err := dbConnection()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	var stmt *sql.Stmt
	var id int

	stmt, err = db.Prepare("insert into cash_collections (tag, sum, status, comment,purpose,create_date, close_date) values ($1,$2,$3,$4,$5,$6,$7) RETURNING id")
	if err != nil {
		return -1, err
	}
	err = stmt.QueryRow(cashCollection.Tag, cashCollection.Sum, cashCollection.Status, cashCollection.Comment, cashCollection.Purpose, cashCollection.CreateDate.Format(time.DateOnly), cashCollection.CloseDate.Format(time.DateOnly)).Scan(&id)

	_ = stmt.Close()
	return id, nil

}

func InfoAboutCashCollection(idCashCollection int) (CashCollection, error) {
	cc := CashCollection{ID: idCashCollection}

	db, err := dbConnection()
	if err != nil {
		return cc, err
	}
	defer db.Close()

	stmt, err := db.Prepare("select tag, sum, status, comment, create_date, close_date, purpose from cash_collections where id =$1")
	if err != nil {

		return cc, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(idCashCollection).Scan(&cc.Tag, &cc.Sum, &cc.Status, &cc.Comment, &cc.CreateDate, &cc.CloseDate, &cc.Purpose)
	return cc, err
}

func UpdateStatusCashCollection(idCashCollection int) error {

	db, err := dbConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("call check_debtors($1)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_ = stmt.QueryRow(idCashCollection)

	return nil
}

type Transaction struct {
	ID               int
	CashCollectionID int
	Sum              float64
	Type             string
	Status           string
	Receipt          string
	MemberID         int64
	Date             time.Time
}

func InfoAboutTransaction(idTransaction int) (Transaction, error) {
	db, err := dbConnection()
	if err != nil {
		return Transaction{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare("select * from transactions where id = $1")
	if err != nil {
		return Transaction{}, err
	}
	defer stmt.Close()

	var t Transaction
	err = stmt.QueryRow(idTransaction).Scan(&t.ID, &t.CashCollectionID, &t.Sum, &t.Type, &t.Status, &t.Receipt, &t.MemberID, &t.Date)

	return t, err
}

func InsertInTransactions(transaction Transaction) (int, error) {
	db, err := dbConnection()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into transactions (cash_collection_id, sum, type, status,receipt, member_id, date) values ($1,$2,$3,$4,$5,$6,$7) RETURNING id")
	if err != nil {
		return -1, err
	}

	var id int
	_ = stmt.QueryRow(transaction.CashCollectionID, transaction.Sum, transaction.Type, transaction.Status, transaction.Receipt, transaction.MemberID, transaction.Date.Format(time.DateOnly)).Scan(&id)
	return id, nil
}

func ChangeStatusTransaction(idTransaction int, status string) error {
	db, err := dbConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("update transactions set status = $1 where id= $2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_ = stmt.QueryRow(status, idTransaction)

	return nil
}

//
//func ExistsFund(tag string) (result bool, err error) {
//	result = false
//
//	db, err := dbConnection()
//	if err != nil {
//		return
//	}
//	defer db.Close()
//
//	stmt, err := db.Prepare("select count(*) from funds where tag=$1")
//	if err != nil {
//		return
//	}
//
//	var count int
//	err = stmt.QueryRow(tag).Scan(&count)
//
//	if (err != nil) || (count == 0) {
//		return
//	}
//
//	result = true
//	return
//}

//
//func CreateDebitingFunds(memberId int64, tag string, sum float64, comment string, purpose string, receipt string) (ok bool, err error) {
//	db, err := dbConnection()
//	if err != nil {
//		return
//	}
//	defer db.Close()
//
//	stmt, err := db.Prepare("select * from  new_deb($1, $2, $3,$4,$5, $6)")
//	if err != nil {
//		return
//	}
//
//	err = stmt.QueryRow(tag, sum, comment, purpose, receipt, memberId).Scan(&ok)
//	return
//}
//
//func GetDebtors(idCashCollection int) (members []int64, err error) {
//	db, err := dbConnection()
//	if err != nil {
//		return
//	}
//	defer db.Close()
//
//	stmt, err := db.Prepare("select member_id from members where member_id not in (select member_id  from transactions where cash_collection_id =$1 and status = 'подтвержден')")
//	if err != nil {
//		return
//	}
//
//	rows, err := stmt.Query(idCashCollection)
//	if err != nil {
//		return
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var member int64
//		if err = rows.Scan(&member); err != nil {
//			return
//		}
//		members = append(members, member)
//	}
//	return
//}

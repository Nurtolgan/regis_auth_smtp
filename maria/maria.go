package maria

import (
	"database/sql"
	"fmt"
	"manager/config"
	"manager/debugger"
	"manager/mongo"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func CreatePlayer(p Person) []byte {
	query := fmt.Sprintf(`insert into person (photo, first_name, last_name, club_id, post, mobile, email) 
	values ('%s', '%s', '%s', %d, '%s', '%s', '%s') returning id;`,
		p.Photo, p.First_name, p.Last_name, p.Club_id, p.Post, p.Mobile, p.Email)
	var id int
	debugger.CheckError("Insert into person", requestRow(query).Scan(&id))
	query = fmt.Sprintf(`insert into player (person_id) values (%d)`, id)
	debugger.CheckError("Insert into player", requestRow(query).Err())
	return []byte("Created")
}

func UpdatePlayer(player Player) {
	query := fmt.Sprintf(`update person set
	first_name='%s',
	last_name='%s',
	mobile='%s',
	email='%s'
	where id=%d`,
		player.First_name, player.Last_name, player.Mobile, player.Email, player.Id)
	row := requestRow(query)
	debugger.CheckError("Update", row.Err())

	query = fmt.Sprintf(`update player set
	gender=%d,
	beach=%d,
	volleyball=%d,
	birth_date='%s',
	birth_place='%s',
	nationality='%s',
	rank='%s',
	height=%d,
	weight=%d,
	spike=%d,
	jump=%d,
	position='%s',
	passport='%s',
	passport_validity='%s',
	first_trainer='%s'
	where person_id=%d`,
		player.Gender, player.Beach, player.Volleyball, player.Birth_date, player.Birth_place,
		player.Nationality, player.Rank, player.Height, player.Weight, player.Spike, player.Jump,
		player.Position, player.Passport, player.Passport_validity, player.First_trainer, player.Id)
	row = requestRow(query)
	debugger.CheckError("Update 2", row.Err())
}

func UpdateTrainer(trainer Person) {
	query := fmt.Sprintf(`update person set
	first_name='%s',
	last_name='%s',
	mobile='%s',
	email='%s'
	where id=%d`,
		trainer.First_name, trainer.Last_name, trainer.Mobile, trainer.Email, trainer.Id)
	row := requestRow(query)
	debugger.CheckError("Update", row.Err())
}

func SignPlayer(cred Credentials) {
	hash, err := bcrypt.GenerateFromPassword([]byte(cred.Password), bcrypt.DefaultCost)
	debugger.CheckError("Generate Password", err)
	query := fmt.Sprintf("insert into credentials (username, password) values ('%s', '%s');", cred.Username, string(hash))
	row := requestRow(query)
	debugger.CheckError("Insert", row.Err())
	debugger.CheckError("Delete from Mongo", mongo.DeleteDocument(cred.Username))
}

func GetCredentials(username string) Credentials {
	query := fmt.Sprintf("select username, password, id, post from `Credentials` where username='%s';", username)
	row := requestRow(query)
	var cred Credentials
	debugger.CheckError("Scan", row.Scan(
		&cred.Username,
		&cred.Password,
		&cred.Person_id,
		&cred.Post))
	return cred
}

func requestRow(query string) *sql.Row {
	mysql := mysql.NewConfig()
	mysql.User = config.Config.Maria_user
	mysql.Passwd = config.Config.Maria_password
	mysql.Net = config.Config.Maria_ip
	mysql.DBName = config.Config.Maria_database
	db, errDB := sql.Open("mysql", mysql.FormatDSN())
	debugger.CheckError("errDB", errDB)
	defer db.Close()

	row := db.QueryRow(query)
	return row
}

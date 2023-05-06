package maria

type Credentials struct {
	Username  string
	Password  string
	Person_id int
	Post      string
}

type Person struct {
	Id         int    `json:"id"`
	Photo      string `json:"photo"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Club_id    int16  `json:"club_id"`
	Club_name  string `json:"club_name"`
	Post       string `json:"post"`
	Mobile     string `json:"mobile"`
	Email      string `json:"email"`
}

type Player struct {
	Id                int    `json:"id"`
	Photo             string `json:"photo"`
	First_name        string `json:"first_name"`
	Last_name         string `json:"last_name"`
	Club_id           int16  `json:"club_id"`
	Club_name         string `json:"club_name"`
	Post              string `json:"post"`
	Mobile            string `json:"mobile"`
	Email             string `json:"email"`
	City              string `json:"city"`
	Gender            int16  `json:"gender"`
	Beach             int16  `json:"beach"`
	Volleyball        int16  `json:"volleyball"`
	Birth_date        string `json:"birth_date"`
	Birth_place       string `json:"birth_place"`
	Nationality       string `json:"nationality"`
	Rank              string `json:"rank"`
	Height            int16  `json:"height"`
	Weight            int16  `json:"weight"`
	Spike             int16  `json:"spike"`
	Jump              int16  `json:"jump"`
	Position          string `json:"position"`
	Passport          string `json:"passport"`
	Passport_validity string `json:"passport_validity"`
	Team_name         string `json:"team_name"`
	First_trainer     string `json:"first_trainer"`
}

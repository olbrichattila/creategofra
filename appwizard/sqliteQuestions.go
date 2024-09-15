package appwizard

var sqliteDBDatabaseQuestion = question{
	key:           "DB_DATABASE",
	prompt:        "Pease provide database path",
	defaultAnswer: "./database/database.sqlite",
	nextQuestion:  &mailQuestion,
}

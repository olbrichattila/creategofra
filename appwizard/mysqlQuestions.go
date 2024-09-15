package appwizard

var mySqlDbHostQuestion = question{
	key:           "DB_HOST",
	prompt:        "Pease provide DB host example: localhost",
	defaultAnswer: "localhost",
	nextQuestion:  &mySqlDbPortQuestion,
}

var mySqlDbPortQuestion = question{
	key:           "DB_PORT",
	prompt:        "Pease provide DB port",
	defaultAnswer: "3306",
	nextQuestion:  &mySqlDbDatabaseNameQuestion,
}

var mySqlDbDatabaseNameQuestion = question{
	key:          "DB_DATABASE",
	prompt:       "Pease provide database name",
	nextQuestion: &mySqlDbUserNameQuestion,
}

var mySqlDbUserNameQuestion = question{
	key:          "DB_USERNAME",
	prompt:       "Pease provide database user name",
	nextQuestion: &mySqlDbPasswordQuestion,
}

var mySqlDbPasswordQuestion = question{
	key:          "DB_PASSWORD",
	prompt:       "Pease provide database password",
	nextQuestion: &mailQuestion,
}

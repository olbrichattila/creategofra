package appwizard

var pgSqlDbHostQuestion = question{
	key:           "DB_HOST",
	prompt:        "Pease provide DB host example: localhost",
	defaultAnswer: "localhost",
	nextQuestion:  &pgSqlDbPortQuestion,
}

var pgSqlDbPortQuestion = question{
	key:           "DB_PORT",
	prompt:        "Pease provide DB port",
	defaultAnswer: "5432",
	nextQuestion:  &pgSqlDbDatabaseNameQuestion,
}

var pgSqlDbDatabaseNameQuestion = question{
	key:          "DB_DATABASE",
	prompt:       "Pease provide database name",
	nextQuestion: &pgSqlDbUserNameQuestion,
}

var pgSqlDbUserNameQuestion = question{
	key:          "DB_USERNAME",
	prompt:       "Pease provide database user name",
	nextQuestion: &pgSqlDbPasswordQuestion,
}

var pgSqlDbPasswordQuestion = question{
	key:          "DB_PASSWORD",
	prompt:       "Pease provide database password",
	nextQuestion: &pgSqlDbSSLModeQuestion,
}

var pgSqlDbSSLModeQuestion = question{
	key:           "DB_SSLMODE",
	prompt:        "Pease provide SSL mode",
	defaultAnswer: "disable",
	nextQuestion:  &mailQuestion,
}

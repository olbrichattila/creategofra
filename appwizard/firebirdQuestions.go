package appwizard

var firebirdDbHostQuestion = question{
	key:           "DB_HOST",
	prompt:        "Pease provide DB host example: localhost",
	defaultAnswer: "localhost",
	nextQuestion:  &firebirdDbPortQuestion,
}

var firebirdDbPortQuestion = question{
	key:           "DB_PORT",
	prompt:        "Pease provide DB port",
	defaultAnswer: "3050",
	nextQuestion:  &firebirdDbDatabaseNameQuestion,
}

var firebirdDbDatabaseNameQuestion = question{
	key:           "DB_DATABASE",
	prompt:        "Pease provide database name",
	defaultAnswer: "/firebird/data/employee.fdb",
	nextQuestion:  &firebirdDbUserNameQuestion,
}

var firebirdDbUserNameQuestion = question{
	key:           "DB_USERNAME",
	prompt:        "Pease provide database user name",
	defaultAnswer: "SYSDBA",
	nextQuestion:  &firebirdDbPasswordQuestion,
}

var firebirdDbPasswordQuestion = question{
	key:           "DB_PASSWORD",
	prompt:        "Pease provide database password",
	defaultAnswer: "masterkey",
	nextQuestion:  &mailQuestion,
}

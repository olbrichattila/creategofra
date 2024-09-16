package appwizard

var databaseQuestions = question{
	key: "DB_CONNECTION",
	prompt: `Pease select database:
  1. MySql
  2. SqLite
  3. PostgresQl
  4. Firebird SQL`,
	answers: answers{
		"1": answer{value: "mysql", nextQuestion: &mySqlDbHostQuestion},
		"2": answer{value: "sqlite", nextQuestion: &sqliteDBDatabaseQuestion},
		"3": answer{value: "pgsql", nextQuestion: &pgSqlDbHostQuestion},
		"4": answer{value: "firebird", nextQuestion: &firebirdDbHostQuestion},
	},
}

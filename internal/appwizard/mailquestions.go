package appwizard

var mailQuestion = question{
	prompt: `Do you want to set up SMTP mail credentials:
  1. Yes
  2. No`,
	answers: answers{
		"1": answer{value: "yes", nextQuestion: &mailUserNameQuestion},
		"2": answer{value: "no"},
	},
}

var mailUserNameQuestion = question{
	key:           "SMTP_USER_NAME",
	prompt:        "Pease provide SMTP user name",
	defaultAnswer: "mailtrap",
	nextQuestion:  &mailPasswordNameQuestion,
}

var mailPasswordNameQuestion = question{
	key:           "SMTP_PASSWORD",
	prompt:        "Please provide SMTP password",
	defaultAnswer: "mailtrap",
	nextQuestion:  &mailHostQuestion,
}

var mailHostQuestion = question{
	key:           "SMTP_HOST",
	prompt:        "Pease provide SMTP host",
	defaultAnswer: "localhost",
	nextQuestion:  &mailPortQuestion,
}

var mailPortQuestion = question{
	key:           "SMTP_PORT",
	prompt:        "Pease provide SMTP host",
	defaultAnswer: "1025",
}

package appwizard

var storageQuestionMap = map[string]*question{
	"file":      nil,
	"db":        nil,
	"redis":     &appRedisHostQuestion,
	"memcached": &appMemcachedHostQuestion,
}

var appRedisHostQuestion = question{
	key:           "REDIS_SERVER_HOST",
	prompt:        "Please provide redis host",
	defaultAnswer: "localhost",
	nextQuestion:  &appRedisPasswordQuestion,
}

var appRedisPasswordQuestion = question{
	key:           "REDIS_PASSWORD",
	prompt:        "Please provide redis password",
	defaultAnswer: "",
	nextQuestion:  &appRedisDbQuestion,
}

var appRedisDbQuestion = question{
	key:           "REDIS_DB",
	prompt:        "Please provide redis DB",
	defaultAnswer: "0",
	nextQuestion:  &appRedisPortQuestion,
}

var appRedisPortQuestion = question{
	key:           "REDIS_PORT",
	prompt:        "Please provide redis port",
	defaultAnswer: "6379",
}

var appMemcachedHostQuestion = question{
	key:           "MEMCACHE_HOST",
	prompt:        "Please provide memcached host",
	defaultAnswer: "localhost",
	nextQuestion:  &appMemcachedPortQuestion,
}

var appMemcachedPortQuestion = question{
	key:           "MEMCACHE_PORT",
	prompt:        "Please provide memcached host",
	defaultAnswer: "11211",
}

var sessionStorageQuestions = question{
	key: "SESSION_STORAGE",
	prompt: `Pease select Session Storage:
  1. File
  2. Redis
  3. Database
  4. Memcached`,
	answers: answers{
		"1": answer{value: "file"},
		"2": answer{value: "redis"},
		"3": answer{value: "db"},
		"4": answer{value: "memcached"},
	},
	nextQuestion: &loggerStorageQuestions,
}

var loggerStorageQuestions = question{
	key: "LOGGER_STORAGE",
	prompt: `Pease select Logger Storage:
  1. File
  2. Redis
  3. Database
  4. Memcached`,
	answers: answers{
		"1": answer{value: "file"},
		"2": answer{value: "redis"},
		"3": answer{value: "db"},
		"4": answer{value: "memcached"},
	},
	nextQuestion: &cacheStorageQuestions,
}

var cacheStorageQuestions = question{
	key: "CACHE_STORAGE",
	prompt: `Pease select Cache Storage:
  1. File
  2. Redis
  3. Database
  4. Memcached`,
	answers: answers{
		"1": answer{value: "file"},
		"2": answer{value: "redis"},
		"3": answer{value: "db"},
		"4": answer{value: "memcached"},
	},
	nextQuestion: &databaseQuestions,
}

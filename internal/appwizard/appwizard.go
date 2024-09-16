package appwizard

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/olbrichattila/creategofra/internal/specio"
)

// EnvData contains answers given by the user
type EnvData struct {
	Key   string
	Value string
}

type answer struct {
	value        string
	nextQuestion *question
}
type answers map[string]answer
type question struct {
	key           string
	prompt        string
	defaultAnswer string
	mandatory     bool
	answers       answers
	nextQuestion  *question
}

func Wizard(envFileName string) []EnvData {
	envContent := getEnvContent(envFileName)

	responses := processQuestions(envContent, appUrlQuestion)

	storages := getStorages(responses)
	for _, storageName := range storages {
		if storageQuestion, ok := storageQuestionMap[storageName]; ok {
			if storageQuestion == nil {
				continue
			}
			responses = append(responses, processQuestions(envContent, *storageQuestion)...)
		}
	}

	envStr := mergeEnv(envContent, responses)
	saveEnvContent(envFileName, envStr)

	return responses
}

func processQuestions(envContent string, q question) []EnvData {
	responses := make([]EnvData, 0)
	currentQuestion := q
	for {
		currentValue := lookupValue(envContent, currentQuestion.key)
		if currentValue == "" {
			currentValue = currentQuestion.defaultAnswer
		}
		answer := selection(currentQuestion, currentValue)
		fmt.Println("")
		if answer == nil {
			break
		}

		if currentQuestion.key != "" && answer.value != "" {
			responses = append(responses, EnvData{Key: currentQuestion.key, Value: answer.value})
		}

		if answer.nextQuestion != nil {
			currentQuestion = *answer.nextQuestion
			continue
		}

		break
	}

	return responses
}

func selection(q question, currentValue string) *answer {
	prompt := ""
	if len(q.answers) > 0 {
		fmt.Println(q.prompt)
		prompt = "Please choose: "
	} else {
		prompt = q.prompt + ": "
	}

	for {
		response := ""
		if len(q.answers) == 0 {
			response = specio.Input(prompt, currentValue)
		} else {
			resolvedAnswer := resolveAnswer(q.answers, currentValue)
			response = specio.Input(prompt, resolvedAnswer)
		}

		if response == "" && q.mandatory {
			if len(q.answers) == 0 {
				fmt.Println("Please provide a value")
			} else {
				fmt.Println("Please select an option")
			}
			continue
		}

		if len(q.answers) > 0 {
			if selected, ok := q.answers[response]; ok {
				if selected.nextQuestion == nil && q.nextQuestion != nil {
					return &answer{value: selected.value, nextQuestion: q.nextQuestion}
				}
				return &selected
			}

			fmt.Println("invalid selection")
			continue
		}

		return &answer{value: response, nextQuestion: q.nextQuestion}
	}
}

func mergeEnv(currentEnv string, data []EnvData) string {
	currentLines := strings.Split(currentEnv, "\n")

	for _, envLine := range data {
		envRow := envLine.Key + "=" + envLine.Value
		if keyId, ok := lookup(currentLines, envLine.Key); ok {
			currentLines[keyId] = envRow
			continue
		}

		currentLines = append(currentLines, envRow)
	}

	return strings.Join(currentLines, "\n")
}

func lookup(lines []string, keyValue string) (int, bool) {
	for i, line := range lines {
		lineParts := strings.Split(line, "=")
		key := strings.TrimSpace(lineParts[0])
		if key == keyValue {
			return i, true
		}
	}

	return 0, false
}

func lookupValue(envContent, keyValue string) string {
	lines := strings.Split(envContent, "\n")
	for _, line := range lines {
		lineParts := strings.Split(line, "=")
		key := strings.TrimSpace(lineParts[0])
		if key == keyValue {
			if len(lineParts) > 1 {
				return strings.TrimSpace(strings.Join(lineParts[1:], "="))
			}
		}
	}

	return ""
}

func getEnvContent(envFileName string) string {
	content, err := os.ReadFile(envFileName)
	if err != nil {
		return ""
	}

	return string(content)
}

func saveEnvContent(envFileName, fileContent string) error {
	return os.WriteFile(envFileName, []byte(fileContent), 0644)
}

func resolveAnswer(a answers, value string) string {
	for key, answer := range a {
		if answer.value == value {
			return key
		}
	}

	return ""
}

func getStorages(data []EnvData) []string {
	re := regexp.MustCompile(`.*_STORAGE`)

	storages := make([]string, 0)
	for _, env := range data {
		if re.MatchString(env.Key) && !sliceContains(storages, env.Value) {
			storages = append(storages, env.Value)
		}
	}

	return storages
}

func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

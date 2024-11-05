package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// Quiz represents the structure of the quiz
type Quiz struct {
	Title     string     `json:"title"`
	Level     string     `json:"level"`
	Time      int        `json:"time"`
	Questions []Question `json:"questions"`
}

// Question represents a single quiz question
type Question struct {
	Question      string   `json:"question"`
	Answers       []string `json:"answers"`
	Points        int      `json:"points"`
	CorrectAnswer string   `json:"correct_answer"`
}

// clearConsole clears the console based on the OS
func clearConsole() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Your platform is not supported to clear the terminal.")
	}
}
func encryptor(){
	//the key should be 32 chars long or string[32]
	key := []byte("32-byte-long-key-for-AES-encrypi")

	// Encrypting the JSON file
	inputFile := "quiz.json"
	encryptedFile := "quiz_encrypted.json"
	err := Encrypt(inputFile, encryptedFile, key)
	if err != nil {
		fmt.Println("Error encrypting file:", err)
		return
	}
	fmt.Println("File encrypted successfully.")

}

// createQuiz allows the user to create a new quiz and saves it to a JSON file
func createQuiz() {
	var quiz Quiz
	fmt.Print("Enter quiz title: ")
	fmt.Scanln(&quiz.Title)

	fmt.Print("Enter quiz level: ")
	fmt.Scanln(&quiz.Level)

	fmt.Print("Enter quiz time (seconds): ")
	fmt.Scanln(&quiz.Time)

	for {
		var question Question
		fmt.Print("Enter question: ")
		fmt.Scanln(&question.Question)

		for i := 0; i < 4; i++ {
			var answer string
			fmt.Printf("Enter answer %d: ", i+1)
			fmt.Scanln(&answer)
			question.Answers = append(question.Answers, answer)
		}

		fmt.Print("Enter points for this question: ")
		fmt.Scanln(&question.Points)

		fmt.Print("Enter the correct answer: ")
		fmt.Scanln(&question.CorrectAnswer)

		quiz.Questions = append(quiz.Questions, question)

		fmt.Print("Add another question? (yes/no): ")
		var moreQuestions string
		fmt.Scanln(&moreQuestions)
		if moreQuestions != "yes" {
			break
		}
	}

	data, err := json.MarshalIndent(quiz, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON data:", err)
		return
	}

	err = ioutil.WriteFile("quiz.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return
	}

	fmt.Println("Quiz saved successfully to quiz.json!")
	//encryptor() ABU - Last Chpont [day before last exams(BE)] - 8/11/2024
}
func main() {
	var mode int
	fmt.Println("Quizzle text-based (terminal/command prompt) interactive question and answer program.")
	fmt.Println("Choose mode:\n\t 1. Execute Quiz\n\t 2. Create Quiz")
	fmt.Scanln(&mode)

	switch mode {
	case 1:
		runQuiz()
	case 2:
		createQuiz()
	default:
		fmt.Println("Invalid mode selected. Exiting...")
	}
}

func runQuiz(){
	fmt.Println("Enter/Paste path to Game File")
	var jsonFilePath string
	fmt.Scanln(&jsonFilePath)
	//decryptor() - ABU [LAST CHECKPOINT] > Tues Jul 9 2024 6:20pm (day before exams)
	// Read the JSON file
	data, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Unmarshal the JSON data into the Quiz struct
	var quiz Quiz
	err = json.Unmarshal(data, &quiz)
	if err != nil {
		fmt.Println("Error parsing JSON data:", err)
		return
	}

	fmt.Println("Welcome to the Quiz Game!")
	fmt.Printf("Title: %s\nLevel: %s\nTime: %d seconds\n", quiz.Title, quiz.Level, quiz.Time)
	fmt.Println("Press Enter to start the quiz...")
	fmt.Scanln()

	score := 0
	//max score feature added by Abu
	maxScore:= 0

	// Iterate over each question
	for i, question := range quiz.Questions {
		clearConsole()
		fmt.Printf("Question %d: %s\n", i+1, question.Question)
		for j, answer := range question.Answers {
			fmt.Printf("%d. %s\n", j+1, answer)
		}
		fmt.Print("Your answer: ")

		var userAnswer int
		fmt.Scan(&userAnswer)

		if userAnswer > 0 && userAnswer <= len(question.Answers) && question.Answers[userAnswer-1] == question.CorrectAnswer {
			fmt.Println("Correct!")
			score += question.Points
		} else {
			fmt.Println("Wrong!")
		}
		maxScore += question.Points

		fmt.Printf("Correct answer: %s\n", question.CorrectAnswer)
		time.Sleep(2 * time.Second)
	}

	clearConsole()
	fmt.Printf("Quiz Over!\nYour Score: %d\nMax Score: %d\n", score, maxScore)
}

//encrypt file
func Encrypt(inputFile, outputFile string, key []byte) error {
	plaintext, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("error creating cipher block: %v", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return fmt.Errorf("error generating IV: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outFile.Close()

	// Encode ciphertext to base64 before writing to file
	encoder := base64.NewEncoder(base64.StdEncoding, outFile)
	defer encoder.Close()

	_, err = encoder.Write(ciphertext)
	if err != nil {
		return fmt.Errorf("error writing encrypted data: %v", err)
	}

	return nil
}

func decryptor(){
	// Decrypting the encrypted JSON file
	key := []byte("32-byte-long-key-for-AES-encrypi")
	encryptedFile := "quiz_encrypted.json"
	decryptedFile := "quiz_decrypted.json"
	err := Decrypt(encryptedFile, decryptedFile, key)
	if err != nil {
		fmt.Println("Error decrypting file:", err)
		return
	}
	fmt.Println("File decrypted successfully.")
}

// Decrypt decrypts a file that was encrypted using AES with a given key
func Decrypt(inputFile, outputFile string, key []byte) error {
	ciphertext, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("error creating cipher block: %v", err)
	}

	if len(ciphertext) < aes.BlockSize {
		return fmt.Errorf("encrypted data is too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outFile.Close()

	_, err = outFile.Write(ciphertext)
	if err != nil {
		return fmt.Errorf("error writing decrypted data: %v", err)
	}

	return nil
}
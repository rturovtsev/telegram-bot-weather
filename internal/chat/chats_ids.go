package chat

import (
	"encoding/json"
	"log"
	"os"
)

var chatIDs []int64

func SaveChatIDs(chatFile string) {
	file, err := os.Create(chatFile)
	if err != nil {
		log.Panicln("Ошибка при сохранении chat ID:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(chatIDs); err != nil {
		log.Panicln("Ошибка при кодировании chat ID:", err)
	}
}

func LoadChatIDs(chatFile string) []int64 {
	file, err := os.Open(chatFile)
	if os.IsNotExist(err) {
		return chatIDs // Файл не существует, ничего страшного
	}
	if err != nil {
		log.Panicln("Ошибка при открытии файла:", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&chatIDs); err != nil {
		log.Panicln("Ошибка при декодировании chat ID:", err)
	}

	return chatIDs
}

func AddChatID(chatID int64, chatFile string) {
	for _, id := range chatIDs {
		if id == chatID {
			return // ChatID уже существует, ничего не делать
		}
	}
	chatIDs = append(chatIDs, chatID)
	SaveChatIDs(chatFile)
}

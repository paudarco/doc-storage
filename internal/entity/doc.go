package entity

import "time"

type Document struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"-" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	IsFile    bool      `json:"file" db:"is_file"`
	Public    bool      `json:"public" db:"public"`
	Mime      string    `json:"mime,omitempty" db:"mime"`
	Grant     []string  `json:"grant,omitempty" db:"grant"`
	CreatedAt time.Time `json:"created" db:"created_at"`

	// Эти поля не хранятся в БД, используются для передачи данных
	JSONData interface{} `json:"json,omitempty" db:"-"` // Для JSON данных
	FileData []byte      `json:"-" db:"-"`              // Для содержимого файла
}

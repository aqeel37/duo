package admin

type User struct {
	Alias1            *string `json:"alias1"`
	Alias2            *string `json:"alias2"`
	Alias3            *string `json:"alias3"`
	Alias4            *string `json:"alias4"`
	Created           int64   `json:"created"`
	Email             string  `json:"email"`
	FirstName         *string `json:"firstname"`
	LastDirectorySync *int64  `json:"last_directory_sync"`
	LastLogin         *int64  `json:"last_login"`
	LastName          *string `json:"lastname"`
	Notes             string  `json:"notes"`
	RealName          *string `json:"realname"`
	Status            string  `json:"status"`
	UserID            string  `json:"user_id"`
	Username          string  `json:"username"`
	IsEnrolled        bool    `json:"is_enrolled"`
}

type MetadateDTO struct {
	NextOffset string `json:"next_offset"`
}

type UserResponseDTO struct {
	Status   string      `json:"stat"`
	Response []User      `json:"response"`
	Metadata MetadateDTO `json:"metadata"`
}

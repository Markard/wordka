package guess

type Request struct {
	Word string `json:"word" validate:"required,len=5"`
}

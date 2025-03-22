package types

type CreateContact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateContactResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Contact struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GetContactsReponse struct {
	Data []*Contact `json:"data"`
}

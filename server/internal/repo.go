package internal

import (
	"sync"

	"github.com/google/uuid"
)

type ContactRepo struct {
	contacts map[string]*Contact
	mux      sync.RWMutex
}

func NewContactRepo() *ContactRepo {
	return &ContactRepo{
		contacts: make(map[string]*Contact),
		mux:      sync.RWMutex{},
	}
}

func (r *ContactRepo) CreateContact(name, email string) (*Contact, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	c := &Contact{
		ID:    uuid.String(),
		Name:  name,
		Email: email,
	}

	r.contacts[c.ID] = c

	return c, nil
}

func (r *ContactRepo) GetContacts() []*Contact {
	r.mux.RLock()
	defer r.mux.RUnlock()

	var contacts []*Contact
	for _, c := range r.contacts {
		contacts = append(contacts, c)
	}

	return contacts
}

package repository

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
	"time"
)

type SessionsRepository struct {
	db db.DB
}

func NewSessionsRepository(db db.DB) SessionsRepository {
	return SessionsRepository{db}
}

func (u *SessionsRepository) ReadSessions() ([]model.Session, error) {
	records, err := u.db.Load("sessions")
	if err != nil {
		return nil, err
	}

	var listSessions []model.Session
	err = json.Unmarshal([]byte(records), &listSessions)
	if err != nil {
		return nil, err
	}

	return listSessions, nil
}

func (u *SessionsRepository) DeleteSessions(tokenTarget string) error {

	// read session
	listSessions, err := u.ReadSessions()
	if err != nil {
		return err
	}

	// Select target token and delete from listSessions
	// TODO: answer here
	for i, session := range listSessions {
		if session.Token == tokenTarget {
			listSessions[i] = model.Session{}
		}
	}

	jsonData, err := json.Marshal(listSessions)
	if err != nil {
		return err
	}

	err = u.db.Save("sessions", jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) AddSessions(session model.Session) error {
	var err error

	// read session
	listSessions, err := u.ReadSessions()
	if err != nil {
		return err
	}

	listSessions = append(listSessions, session)
	// struct menjadi json
	jsonData, err := json.Marshal(listSessions)
	if err != nil {
		panic(err)
	}

	// menulis data ke file
	err = u.db.Save("sessions", jsonData)
	if err != nil {
		return err
	}
	return err // TODO: replace this
}

func (u *SessionsRepository) CheckExpireToken(token string) (model.Session, error) {

	session, err := u.TokenExist(token)
	if err != nil {
		err = fmt.Errorf("token doesn't exist")
		return model.Session{}, err
	}

	if u.TokenExpired(session) {
		// expired < now

		// delete session di listsession
		err := u.DeleteSessions(token)
		if err != nil {
			err = fmt.Errorf("error when deleting session")
			return session, err
		}

		return model.Session{}, fmt.Errorf("Token is Expired!")

	} else {
		// expired > noew
		return session, nil
	}

	// //find session
	// for i, session := range listSessions {
	// 	if session.Token == token {
	// 		if u.TokenExpired(session) {
	// 			return session, nil
	// 		} else {
	// 			// TODO: replace this

	// 		}
	// 	}
	// }
	// // struct menjadi json
	// jsonData, err := json.Marshal(listSessions)
	// if err != nil {
	// 	return model.Session{}, err
	// }

	// // menulis data ke file
	// err = u.db.Save("sessions", jsonData)
	// if err != nil {
	// 	return model.Session{}, err
	// }

	// return model.Session{}, nil // TODO: replace this
}

func (u *SessionsRepository) ResetSessions() error {
	err := u.db.Reset("sessions", []byte("[]"))
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) TokenExist(req string) (model.Session, error) {
	listSessions, err := u.ReadSessions()
	if err != nil {
		return model.Session{}, err
	}
	for _, element := range listSessions {
		if element.Token == req {
			return element, nil
		}
	}
	return model.Session{}, fmt.Errorf("Token Not Found!")
}

func (u *SessionsRepository) TokenExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}

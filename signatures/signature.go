package signatures

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

type Signature struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (s Signature) WriteToFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return err
	}
	enc := json.NewEncoder(file)
	return enc.Encode(&s)
}

func SaveOverFile(signList []Signature, filename string) error {
	if err := os.Remove(filename); err != nil {
		return err
	}
	for _, sign := range signList {
		if err := sign.WriteToFile(filename); err != nil {
			return err
		}
	}
	return nil
}

func NewSignature(name string) Signature {
	return Signature{
		Name:      name,
		CreatedAt: time.Now(),
	}
}

func WriteNewSignature(name, filename string) (Signature, error) {
	s := NewSignature(name)
	return s, s.WriteToFile(filename)
}

func ReadFromFile(filename string) ([]Signature, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModeAppend)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(file)
	signatures := []Signature{}
	for {
		s := Signature{}
		err := dec.Decode(&s)
		//Don't want to have any bad members
		if err != nil {
			if err == io.EOF {
				return signatures, nil
			} else {
				return nil, err
			}
		}
		signatures = append(signatures, s)
	}
	return signatures, nil
}

package signatures

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

//type to encapsulate a signature object
type Signature struct {
	//The "signature" string
	Name string `json:"name"`
	//The point in time at which this signature was created
	CreatedAt time.Time `json:"created_at"`
}

//Takes in a signature, and the name of a file, and appends this signature object to the end of the file
//Returns any errors with opening the file or encoding the object
func (s Signature) WriteToFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return err
	}
	enc := json.NewEncoder(file)
	return enc.Encode(&s)
}

//Takes an array of signatures, and the name of a file
//and saves this array of objects by overwriting the file
//Returns any errors with replacing the file or encoding the object
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

//Returns a new signature object with the given name and the current time
func NewSignature(name string) Signature {
	return Signature{
		Name:      name,
		CreatedAt: time.Now(),
	}
}

//Takes a name and a filename, and immediately creates a signature and writes it to the file
//Returns the new signature, and any errors with writing it to the file
func WriteNewSignature(name, filename string) (Signature, error) {
	s := NewSignature(name)
	return s, s.WriteToFile(filename)
}

//Takes a filename of a file with signature objects encoded in it, and returns
//the slice of signature objects written to the file
//Returns an error if there is an issue decoding the file
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
		//Quit early if there are any badly formatted members
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

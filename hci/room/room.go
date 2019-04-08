package room

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Amount struct {
	UserId int `json:"user_id"`
	Amount int `json:"amount"`
}

type Item struct {
	Id      int      `json:"id"`
	Price   float64  `json:"price"`
	Amounts []Amount `json:"amount"`
}

type User struct {
	UserId int `json:"user_id"`
}

type Bill struct {
	UserId int     `json:"user_id"`
	Bill   float64 `json:"bill"`
}

type Room struct {
	Code      string  `json:"code"`
	Tax       float64 `json:"tax"`
	Tip       float64 `json:"tip"`
	Total     float64 `json:"total"`
	UserCount int     `json:"user_count"`
	Items     []Item  `json:"items"`
}

type Request struct {
	Code   string  `json:"room_code"`
	UserId int     `json:"user_id"`
	ItemId int     `json:"item_id"`
	Tax    float64 `json:"tax"`
	Tip    float64 `json:"tip"`
	Total  float64 `json:"total"`
	Price  float64 `json:"price"`
	Amount int     `json:"amount"`
}

var Rooms []Room

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randCode() string {
	b := make([]rune, 4)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func createRoom(tip, tax, total float64) *Room {
	conflictFound := true
	newCode := randCode()
	for conflictFound {
		conflictFound = false
		for _, room := range Rooms {
			if room.Code == newCode {
				newCode = randCode()
				conflictFound = true
				break
			}
		}
	}

	newRoom := Room{
		Code:  newCode,
		Tip:   tip,
		Tax:   tax,
		Total: total,
	}

	Rooms = append(Rooms, newRoom)

	return &newRoom
}

func getRoom(roomCode string) (*Room, error) {
	for i, r := range Rooms {
		if r.Code == roomCode {
			return &Rooms[i], nil
		}
	}
	return nil, errors.New("No such room code")
}

func (r *Room) addUser() (*User, error) {
	newUser := User{
		UserId: r.UserCount,
	}
	r.UserCount = r.UserCount + 1
	fmt.Println(r)
	return &newUser, nil
}

func (r *Room) addItem(item_id, user_id, amount int, price float64) (*Item, error) {
	newItem := Item{
		Id:      len(r.Items),
		Price:   price,
		Amounts: []Amount{Amount{UserId: user_id, Amount: amount}},
	}
	for _, i := range r.Items {
		if i.Id == item_id {
			i.Amounts = append(i.Amounts, Amount{UserId: user_id, Amount: amount})
			newItem = i
		}
	}
	r.Items = append(r.Items, newItem)
	return &newItem, nil
}

func (r *Room) removeParticipation(user_id int, item_id int) (*Item, error) {
	for ind, item := range r.Items {
		if item.Id == item_id {
			newAmounts := []Amount{}
			for _, a := range item.Amounts {
				if a.UserId != user_id {
					newAmounts = append(newAmounts, a)
				}
			}
			r.Items[ind].Amounts = newAmounts
		}
		return &r.Items[ind], nil
	}
	return nil, errors.New("Item ID not found")
}

func (r *Room) calculateBill(user_id int) (*Bill, error) {
	return &Bill{UserId: user_id, Bill: r.Total / float64(r.UserCount)}, nil
}

func CreateRoom(c *gin.Context) {
	var requestDecoded Request
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&requestDecoded); err != nil {
		fmt.Println(requestDecoded)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r := createRoom(requestDecoded.Tip, requestDecoded.Tax, requestDecoded.Total)
	c.JSON(http.StatusOK, r)
}

func CreateItem(c *gin.Context) {
	var requestDecoded Request
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&requestDecoded); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r, err := getRoom(requestDecoded.Code)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	i, err := r.addItem(requestDecoded.ItemId, requestDecoded.UserId, requestDecoded.Amount, requestDecoded.Price)

	if err == nil {
		c.JSON(http.StatusOK, i)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func GetRoom(c *gin.Context) {
	var requestDecoded Request
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&requestDecoded); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r, err := getRoom(requestDecoded.Code)

	if err == nil {
		c.JSON(http.StatusOK, r)
	} else {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

func RemoveParticipation(c *gin.Context) {
	var requestDecoded Request
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&requestDecoded); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r, err := getRoom(requestDecoded.Code)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	i, err := r.removeParticipation(requestDecoded.UserId, requestDecoded.ItemId)

	if err == nil {
		c.JSON(http.StatusOK, i)
	} else {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

func DeleteRooms(c *gin.Context) {
	Rooms = []Room{}
	c.AbortWithStatus(http.StatusOK)
}

func JoinRoom(c *gin.Context) {
	var requestDecoded Request
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&requestDecoded); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r, err := getRoom(requestDecoded.Code)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	u, err := r.addUser()

	if err == nil {
		c.JSON(http.StatusOK, u)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func GetBill(c *gin.Context) {
	var requestDecoded Request
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&requestDecoded); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r, err := getRoom(requestDecoded.Code)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	b, err := r.calculateBill(requestDecoded.UserId)

	if err == nil {
		c.JSON(http.StatusOK, b)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

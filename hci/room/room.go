package room

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Amount struct {
	UserId int `json:"user_id"`
	Amount int `json:"amount"`
}

type Item struct {
	Id      int      `json:"id"`
	Image   string   `json:"image"`
	Price   float64  `json:"price"`
	Amounts []Amount `json:"amount"`
}

type User struct {
	UserId int `json:"user_id"`
}

type Bill struct {
	UserId int     `json:"user_id"`
	Bill   float64 `json:"bill"`
	Items  []Item  `json:"items"`
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
	Image  string  `json:"image"`
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

func (r *Room) addItem(image string) (*Item, error) {
	newItem := Item{
		Id:    len(r.Items),
		Image: image,
	}
	r.Items = append(r.Items, newItem)
	return &newItem, nil
}

func (r *Room) editItem(item_id int, price float64) (*Item, error) {
	for ind, i := range r.Items {
		if i.Id == item_id {
			r.Items[ind].Price = price
			return &r.Items[ind], nil
		}
	}
	return nil, errors.New("Item ID not found")
}

func (r *Room) addParticipation(user_id, item_id, amount int) (*Item, error) {
	for ind, item := range r.Items {
		if item.Id == item_id {
			r.Items[ind].Amounts = append(item.Amounts, Amount{UserId: user_id, Amount: amount})
			return &r.Items[ind], nil
		}
	}
	return nil, errors.New("Item ID not found")
}

func (r *Room) editParticipation(user_id, item_id, amount int) (*Item, error) {
	for ind, item := range r.Items {
		if item.Id == item_id {
			for ind2, a := range item.Amounts {
				if a.UserId == user_id {
					r.Items[ind].Amounts[ind2] = Amount{UserId: user_id, Amount: amount}
					return &r.Items[ind], nil
				}
			}
		}
	}
	return nil, errors.New("Item ID or User ID not found")
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
	return nil, errors.New("Item ID or User ID not found")
}

func (r *Room) calculateBill(user_id int) (*Bill, error) {
	participatedItems := []Item{}
	for _, i := range r.Items {
		for _, a := range i.Amounts {
			if a.UserId == user_id {
				participatedItems = append(participatedItems, i)
				break
			}
		}
	}

	return &Bill{UserId: user_id, Bill: r.Total / float64(r.UserCount), Items: participatedItems}, nil
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

	i, err := r.addItem(requestDecoded.Image)

	if err == nil {
		c.JSON(http.StatusOK, i)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func EditItem(c *gin.Context) {
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

	i, err := r.editItem(requestDecoded.ItemId, requestDecoded.Price)

	if err == nil {
		c.JSON(http.StatusOK, i)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func GetRoom(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r, err := getRoom(code)

	if err == nil {
		c.JSON(http.StatusOK, r)
	} else {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

func AddParticipation(c *gin.Context) {
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

	i, err := r.addParticipation(requestDecoded.UserId, requestDecoded.ItemId, requestDecoded.Amount)

	if err == nil {
		c.JSON(http.StatusOK, i)
	} else {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

func EditParticipation(c *gin.Context) {
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

	i, err := r.editParticipation(requestDecoded.UserId, requestDecoded.ItemId, requestDecoded.Amount)

	if err == nil {
		c.JSON(http.StatusOK, i)
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
	code := c.Param("code")
	if code == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userIdString := c.Param("user")
	if userIdString == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r, err := getRoom(code)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	b, err := r.calculateBill(userId)

	if err == nil {
		c.JSON(http.StatusOK, b)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

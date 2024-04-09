package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"
	// "github.com/mattn/go-sqlite3"
	resp "github.com/nicklaw5/go-respond"
)

// Requests

type (
	Username struct {
		string
	}
	UserId struct {
		int
	}

	UserData struct {
		Username *string `json:"username"`
		UserId   *int    `json:"user_id"`
	}

	createUser struct {
		Username string `json:"username" validate:"required"`
	}
)

// Handlers

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
}

type CreateUserResponse struct {
	Id int `json:"user_id"`
}

func HandleCreateUser(rw http.ResponseWriter, r *http.Request) {
	body, err := ParseValidate[CreateUserRequest](r)
	if err != nil {
		WriteParseValidateErr(rw, err)
		return
	}

	db, err := NewUserStore()
	if err != nil {
		Bail(err, rw)
		return
	}

	id, err := db.CreateUser(NewUser{
		Username: body.Username,
	})
	if err != nil {
		RespE(err, rw)
		return
	}

	resp.NewResponse(rw).Accepted(CreateUserResponse{Id: id})
}

type GetUserRequest struct {
	Id int `json:"user_id" validate:"required"`
}

func HandleGetUser(rw http.ResponseWriter, r *http.Request) {
	body, err := ParseValidate[GetUserRequest](r)
	if err != nil {
		WriteParseValidateErr(rw, err)
		return
	}

	db, err := NewUserStore()
	if err != nil {
		log.Error("Could not set up user store for request", "err", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := db.RetrieveUser(body.Id)
	if err != nil {
		resp.NewResponse(rw).InternalServerError(MsgE(err))
		return
	}

	resp.NewResponse(rw).Ok(user)
}

type DeleteUserRequest struct {
	Id int `json:"user_id" validate:"required"`
}

func HandleDeleteUser(rw http.ResponseWriter, r *http.Request) {
	body, err := ParseValidate[DeleteUserRequest](r)
	if err != nil {
		WriteParseValidateErr(rw, err)
		return
	}
	db, err := NewUserStore()
	if err != nil {
		Bail(err, rw)
		return
	}
	// Validate user exists
	var user struct {
		// TODO: this should contain all available fields of user
		Id string
	}
	if err := db.crud.retrieve.QueryRow(body.Id).Scan(&user); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err = db.crud.delete.Exec(body.Id); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Return with user as a last chance for client to use data
	resp.NewResponse(rw).Ok(user)
}

func RespE(err error, rw http.ResponseWriter) {
	// Handles errors then bails to internal server err
	response := resp.NewResponse(rw)
	switch err {
	case UserExistsError:
		response.BadRequest(MsgE(UserExistsError))
	default:
		Bail(err, rw)
	}
}

// Response builders

type (
	MessageResponse struct {
		Message string `json:"message"`
	}
)

func Msg(message string) MessageResponse {
	return MessageResponse{
		Message: message,
	}
}

func MsgE(err error) MessageResponse {
	return MessageResponse{
		Message: err.Error(),
	}
}

func Bail(err error, rw http.ResponseWriter) {
	resp.NewResponse(rw).DefaultMessage().InternalServerError(nil)
	log.Error(err.Error())
}

func ParseValidate[V interface{}](r *http.Request) (*V, error) {
	var val V
	if err := json.NewDecoder(r.Body).Decode(&val); err != nil {
		return nil, err
	}

	v := validator.New()
	if err := v.Struct(val); err != nil {
		return nil, err
	}

	return &val, nil
}

type Dict map[string]interface{}
type List []interface{}

func WriteParseValidateErr(rw http.ResponseWriter, err error) {
	if vErr, ok := err.(*validator.InvalidValidationError); ok {
		WriteValidationErr(rw, vErr)
		return
	}
	resp.NewResponse(rw).BadRequest(DictErr(
		"could not parse",
		err.Error(),
	))
}

func DictErr(e string, details ...interface{}) Dict {
	return Dict{
		"details": details,
		"error":   e,
	}
}

func WriteValidationErr(rw http.ResponseWriter, err error) {
	errLines := make([]string, 0)
	errs := err.(validator.ValidationErrors)

	for _, fieldErr := range errs {
		l := fmt.Sprintf("field %s: %s", fieldErr.Field(), fieldErr.Tag())
		errLines = append(errLines, l)
	}

	msg := map[string]interface{}{
		"error":   "could not validate request body",
		"details": errLines,
	}

	resp.NewResponse(rw).BadRequest(msg)
}

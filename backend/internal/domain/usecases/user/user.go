package user

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/Louffty/green-code-moscow/internal/domain/dto"
	"github.com/Louffty/green-code-moscow/internal/domain/entities"
	"io/ioutil"
	"log"
	"net/http"
)

type Service interface {
	GetByUUID(ctx context.Context, uuid string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) (*entities.User, error)
}

type EventsUserUseCase interface {
	GetAllByUser(ctx context.Context, userUUID string, limit, offset int) (*dto.ReturnUsersEvents, error)
}

type userUseCase struct {
	userService       Service
	eventsUserUseCase EventsUserUseCase
}

func NewUserUseCase(userService Service, eventsUserUseCase EventsUserUseCase) *userUseCase {
	return &userUseCase{
		userService:       userService,
		eventsUserUseCase: eventsUserUseCase,
	}
}

func (u *userUseCase) GetByUUID(ctx context.Context, uuid string) (*dto.ReturnUser, error) {
	user, err := u.userService.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	events, errGetAllEvents := u.eventsUserUseCase.GetAllByUser(ctx, uuid, 1000, 0)
	if errGetAllEvents != nil {
		return nil, err
	}

	return &dto.ReturnUser{
		UUID:        user.UUID,
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
		CoinsAmount: user.CoinsAmount,
		Rate:        user.Rate,
		Events:      *events,
		IsVerified:  user.IsVerified,
	}, nil
}

func (u *userUseCase) VerifiedUser(ctx context.Context, uuid string, organisation string) (*entities.User, error) {
	user, err := u.userService.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	url := "https://reestrs.minjust.gov.ru/rest/registry/ac648356-fe24-e11a-ceb0-87718bb81ed4/values"

	requestBody := map[string]interface{}{
		"offset": 0,
		"limit":  20,
		"search": organisation,
		"facets": map[string]interface{}{},
		"sort":   []interface{}{},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "insomnia/10.0.0")
	req.Header.Set("Origin", "https://www.minjust.gov.ru")
	req.Header.Set("Sec-Fetch-Mode", "cors")

	// Create an HTTP client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Disable SSL certificate verification
		},
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	type Value struct {
		GridS        string `json:"grid_s"`
		Field9S      string `json:"field_9_s"`
		Field6S      string `json:"field_6_s"`
		Field5S      string `json:"field_5_s"`
		Field8S      string `json:"field_8_s"`
		Field7S      string `json:"field_7_s"`
		Field2S      string `json:"field_2_s"`
		Field1S      string `json:"field_1_s"`
		Field4S      string `json:"field_4_s"`
		Field3S      string `json:"field_3_s"`
		Version      int64  `json:"_version_"`
		Field5DT     string `json:"field_5_dt"`
		ID           string `json:"id"`
		LastModified int64  `json:"lastModified_l"`
	}

	type Count struct {
		Count int    `json:"count"`
		Name  string `json:"name"`
	}

	type Facet struct {
		Name   string  `json:"name"`
		Counts []Count `json:"counts"`
	}

	type Response struct {
		Offset int     `json:"offset"`
		Limit  int     `json:"limit"`
		Size   int     `json:"size"`
		Values []Value `json:"values"`
		Facets []Facet `json:"facets"`
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Unmarshal the JSON response
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON response: %v", err)
	}

	// Check if there is at least one object in the `values` array
	if len(response.Values) > 0 {
		user.IsVerified = true
		return u.userService.Update(ctx, user)
	}

	return user, nil
}

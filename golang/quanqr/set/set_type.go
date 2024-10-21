package set_qr

import (
    "time"
)

type SetSnapshot struct {
    ID          int32       `json:"id"`
    Name        string      `json:"name"`
    Description string      `json:"description"`
    Dishes      []SetDish   `json:"dishes"`
    UserID      *int32      `json:"userId,omitempty"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
    SetID       int32       `json:"set_id"`
    IsPublic    bool        `json:"is_public"`
    Image       string      `json:"image"`
}

type Set struct {
    ID          int64       `json:"id"`
    Name        string      `json:"name"`
    Description string      `json:"description"`
    Dishes      []SetDish   `json:"dishes"`
    UserID      *int32      `json:"userId,omitempty"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
    IsFavourite bool        `json:"is_favourite"`
    LikeBy      []int64     `json:"like_by"`
    IsPublic    bool        `json:"is_public"`
    Image       string      `json:"image"`
}

type SetDish struct {
    DishID   int64 `json:"dish_id"`
    Quantity int64 `json:"quantity"`
}

type CreateSetRequest struct {
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Dishes      []SetDish `json:"dishes"`
    UserID      int32     `json:"userId"`
    IsPublic    bool      `json:"is_public"`
    Image       string    `json:"image"`
}

type UpdateSetRequest struct {
    ID          int32     `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Dishes      []SetDish `json:"dishes"`
    IsPublic    bool      `json:"is_public"`
    Image       string    `json:"image"`
}

type SetResponse struct {
    Data Set `json:"data"`
}

type SetListResponse struct {
    Data []Set `json:"data"`
}

type SetIDParam struct {
    ID int32 `json:"id"`
}
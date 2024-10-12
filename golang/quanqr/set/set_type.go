package set_qr

import "time"

type Dish struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Price       int       `json:"price"`
    Description string    `json:"description"`
    Image       string    `json:"image"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type SetSnapshot struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Dishes      []SetDish `json:"dishes"`
    UserID      *int      `json:"userId,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    SetID       int       `json:"set_id"`
    IsPublic    bool      `json:"is_public"` // add
}

type Set struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Dishes      []SetDish `json:"dishes"`
    UserID      *int32    `json:"userId,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    IsFavourite bool      `json:"is_favourite"`
    LikeBy      []int64   `json:"like_by"`
    IsPublic    bool      `json:"is_public"` // add
}

type SetDish struct {
    DishID   int64 `json:"dish_id"`  // Changed to only store the dish id
    Quantity int   `json:"quantity"`
}

type CreateSetRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Dishes      []struct {
        Dish     Dish `json:"dish"`
        Quantity int  `json:"quantity"`
    } `json:"dishes"`
    UserID   *int32 `json:"userId,omitempty"`
    IsPublic bool   `json:"is_public"` // add 
}




type UpdateSetRequest struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Dishes      []struct {
        Dish     Dish `json:"dish"`
        Quantity int  `json:"quantity"`
    } `json:"dishes"`
    IsPublic bool `json:"is_public"` // add
}
type SetResponse struct {
    Data    Set    `json:"data"`

}

type SetListResponse struct {
    Data    []Set  `json:"data"`

}







package user

// import "english-ai-full/ecomm-grpc/models"
// type CreateUserRequest struct {
//     Username            string `validate:"omitempty,min=2,max=100" json:"name"`
//     Email    string `validate:"required,min=2,max=100" json:"email"`
//     Password string `validate:"required,min=2,max=100" json:"password"`
//     Image           string

   
//     PhoneNumber     string
//     StreetAddress string

  
// }

// type UpdateUserRequest struct {
//     ID       int64    `validate:"required"`
//     Username     string `validate:"required,max=200,min=2" json:"name"`
//     Email    string `validate:"required,min=2,max=100" json:"email"`
//     HashedPassword string `validate:"required,min=2,max=100" json:"password"`
//     Image           string
//     FavoriteIds     string // empty for when createing user
//     PhoneNumber     string
//     StreetAddress string
//     Orders    []models.Order



//     // Add other fields as needed
// }


type LoginRequest struct {
    Email    string `validate:"required,max=200,min=2" json:"email"`
    Password string `validate:"required,min=2,max=100" json:"password"`
}

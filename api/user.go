package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	db "github.com/Gokul-B12/money-txn/db/sqlc"
	"github.com/Gokul-B12/money-txn/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type UserResponse struct {
	Username string `json:"username"`
	//HashedPassword    string    `json:"hashed_password"`  //we should not return password to user.. not a good sec practice
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	password, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: password,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	//arg = db.CreateUserParams{} //testing purpose

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			//log.Println(pqErr.Code.Name())

			switch pqErr.Code.Name() {
			case "unique_violation": //why we are using this is .. we have unique constraint on username and email column on users table.
				log.Println(pqErr.Code.Name())
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := userResponse(user)

	ctx.JSON(http.StatusOK, res)

}

func userResponse(user db.User) UserResponse {
	return UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	User        UserResponse
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	//if user exists, we are checking the passwords to make sure he entered correct password.
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	accessToken, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userRes := loginUserResponse{
		AccessToken: accessToken,
		User:        userResponse(user),
	}

	ctx.JSON(http.StatusOK, userRes)
}

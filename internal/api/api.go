package api

// @title Bank API
// @version 1.0
// @description This is a sample service for managing bank users and accounts
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-jwt-token
import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "go-bank-v2/docs"
	. "go-bank-v2/internal/types"
	"log"
	"net/http"
	"strconv"
)

func (s *Server) Run() {

	router := gin.Default()

	router.POST("/login", s.handleLogin)
	router.POST("/users", s.handleCreateUser)
	router.GET("/users/:id", withJWTAuth(s.handleGetUser, s.store, false))
	router.GET("/users/:id/accounts", withJWTAuth(s.handleGetAllUserAccounts, s.store, false))
	router.DELETE("/users/:id", withJWTAuth(s.handleDeleteUser, s.store, true))
	router.GET("/accounts/:accId", withJWTAuth(s.handleGetAccount, s.store, false))
	router.DELETE("/accounts/:accId", withJWTAuth(s.handleDeleteAccount, s.store, true))
	router.GET("/accounts", withJWTAuth(s.handleGetAllAccounts, s.store, true))
	router.POST("/accounts", withJWTAuth(s.handleCreateAccount, s.store, false))
	router.PATCH("/accounts/:accId", withJWTAuth(s.handleUpdateAccount, s.store, true))
	router.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	fmt.Println("JSON API server running on port:", s.listenAddr)
	err := router.Run(s.listenAddr)
	if err != nil {
		return
	}
}

type Error struct {
	Error string `json:"error"`
}

// @Summary Get All User Accounts
// @Description Fetch all accounts for a user by their ID
// @Tags account
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} types.Account
// @Failure 400 {object} Error "Bad Request"
// @Failure 401 {object} Error "Unauthorized"
// @Failure 403 {object} Error "Forbidden"
// @Failure 500 {object} Error "Internal Server Error"
// @Router /users/{id}/accounts [get]
func (s *Server) handleGetAllUserAccounts(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Error{Error: "Invalid ID"})
		return
	}

	accounts, err := s.store.GetAccounts(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)

}

// @Summary Get User by ID
// @Description Fetch a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} types.UserDto
// @Failure 400 {object} Error "Bad Request"
// @Failure 401 {object} Error "Unauthorized"
// @Failure 404 {object} Error "Not Found"
// @Failure 500 {object} Error "Internal Server Error"
// @Router /users/{id} [get]
func (s *Server) handleGetUser(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.JSON(http.StatusMethodNotAllowed, Error{Error: "Method not allowed"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Error{Error: "Invalid ID"})
		return
	}

	user, err := s.store.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, Error{Error: "Invalid ID"})
		return
	}

	userDto := UserDto{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		MemberSince: user.MemberSince,
		Role:        user.Role,
	}
	c.JSON(http.StatusOK, userDto)
}

// @Security ApiKeyAuth
// @Summary Create an Account
// @Description Create an account for user with given ID
// @Tags account
// @Accept json
// @Produce json
// @Query id path string true "User ID"
// @Success 200 {object} types.Account
// @Failure 400 {object} Error "Bad Request"
// @Router /accounts [post]
func (s *Server) handleCreateAccount(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, Error{Error: "Method not allowed"})
		return
	}
	id, err := strconv.Atoi(c.Query("ownerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Error{Error: "Couldn't create an account!"})
		fmt.Println(err)
		return
	}

	acc := NewAccount(id)
	err = s.store.CreateAccount(acc)
	if err != nil {
		c.JSON(http.StatusMethodNotAllowed, Error{Error: "Couldn't create account!"})
		log.Fatal(err)
		return
	}
	c.JSON(http.StatusOK, acc)
}

// @Summary Login
// @Description Used to log in a user
// @Tags users
// @Accept json
// @Produce json
// @Param id query string true "User ID"
// @Param password query string true "User password"
// @Success 200 {object} types.LoginResponse
// @Failure 400 {object} Error "Bad Request"
// @Router /login [post]
func (s *Server) handleLogin(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, Error{Error: "Method not allowed"})
		return
	}
	id, err := strconv.Atoi(c.Query("id"))
	password := c.Query("password")

	if err != nil || password == "" {
		c.JSON(http.StatusUnauthorized, Error{Error: "Unauthorized"})
		return
	}
	req := LoginRequest{
		ID:       id,
		Password: password}

	user, err := s.store.GetUserByID(req.ID)
	if err != nil {
		c.JSON(http.StatusMethodNotAllowed, Error{Error: "Method not allowed"})
		return
	}

	if !user.ValidPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, Error{Error: "Unauthorized"})
		return
	}

	token, err := createJWT(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, Error{Error: "Something went wrong!"})
		log.Println(err)
		return
	}

	resp := LoginResponse{
		Token: token,
		ID:    user.ID,
	}

	c.JSON(http.StatusOK, resp)

}

// @Summary Create user
// @Description Creates a new user
// @Tags users
// @Accept json
// @Produce json
// @Param id query string true "User ID"
// @Param password query string true "User password"
// @Success 200 {object} types.User
// @Failure 400 {object} Error "Bad Request"
// @Router /users [post]
func (s *Server) handleCreateUser(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, Error{Error: "Error while creating a user"})
		return
	}
	firstName := c.Query("firstName")
	lastName := c.Query("lastName")
	password := c.Query("password")
	role := c.Query("role")
	if firstName == "" || lastName == "" || password == "" {
		c.JSON(http.StatusMethodNotAllowed, Error{Error: "Error while creating user"})
		return
	}
	req := CreateUserRequest{
		FirstName: firstName,
		LastName:  lastName,
		Password:  password,
		Role:      role,
	}

	user, err := NewUser(req.FirstName, req.LastName, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{Error: "Error while creating a user"})
		log.Println(err)
		return
	}
	if err := s.store.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, Error{Error: "Error while creating a user"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Security ApiKeyAuth
// @Summary Get Account
// @Description Fetch an account by account number
// @Tags account
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} types.Account
// @Failure 400 {object} Error "Bad Request"
// @Router /accounts/{id} [get]
func (s *Server) handleGetAccount(c *gin.Context) {
	accNumStr := c.Param("accId")
	accNum, err := strconv.Atoi(accNumStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Error{Error: "Invalid Account Number"})
		return
	}
	account, err := s.store.GetAccountByNumber(accNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

// @Summary Update Account
// @Description Update an account's balance by account number
// @Tags account
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Param newBalance query number true "New Balance"
// @Security ApiKeyAuth
// @Success 200 {object} types.Account
// @Failure 400 {object} Error "Bad Request"
// @Failure 401 {object} Error "Unauthorized"
// @Failure 403 {object} Error "Forbidden"
// @Failure 404 {object} Error "Not Found"
// @Failure 500 {object} Error "Internal Server Error"
// @Router /accounts/{id} [patch]
func (s *Server) handleUpdateAccount(c *gin.Context) {
	accNumStr := c.Param("accId")
	accNum, err := strconv.Atoi(accNumStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Error{Error: "Invalid Account Number"})
		return
	}
	account, err := s.store.GetAccountByNumber(accNum)

	if err != nil || account == nil {
		c.JSON(http.StatusNotFound, Error{Error: "No such account"})
		return
	}

	newBalance, err := strconv.ParseFloat(c.Query("newBalance"), 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{Error: "Invalid balance update!"})
		return
	}
	err = s.store.UpdateAccountBalance(account, newBalance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	account.Balance = newBalance
	c.JSON(http.StatusOK, account)
}

// @Summary Delete Account
// @Description Delete an account by account number
// @Tags account
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Security ApiKeyAuth
// @Success 200 {object} string
// @Failure 400 {object} Error "Bad Request"
// @Failure 401 {object} Error "Unauthorized"
// @Failure 403 {object} Error "Forbidden"
// @Failure 404 {object} Error "Not Found"
// @Failure 500 {object} Error "Internal Server Error"
// @Router /accounts/{id} [delete]
func (s *Server) handleDeleteAccount(c *gin.Context) {
	accNumStr := c.Param("accId")
	accNum, err := strconv.Atoi(accNumStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Error{Error: "Invalid Account Number"})
		return
	}
	account, err := s.store.GetAccountByNumber(accNum)

	if err != nil || account == nil {
		c.JSON(http.StatusNotFound, Error{Error: "No such account"})
		return
	}

	err = s.store.DeleteAccount(accNum)
	if err != nil {
		c.JSON(http.StatusAccepted, Error{Error: "Scheduled for deletion!"})
		return
	}

	c.JSON(http.StatusOK, "Account successfully deleted!")
}

// @Security ApiKeyAuth
// @Summary Delete User by ID
// @Description Deletes a user by their ID, and deletes all the users accounts
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} types.Account
// @Failure 400 {object} Error "Bad Request"
// @Failure 403 {object} Error "Forbidden"
// @Router /users/{id} [delete]
func (s *Server) handleDeleteUser(c *gin.Context) {
	if c.Request.Method != "DELETE" {
		c.JSON(http.StatusMethodNotAllowed, Error{Error: "Method not allowed"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Error{Error: "Invalid ID"})
		return
	}

	accounts, err := s.store.GetAccounts(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{Error: "Can't get accounts"})
		return
	}

	for _, acc := range accounts {
		err = s.store.DeleteAccount(acc.AccountNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Error{Error: "Can't delete accounts"})
			return
		}
	}

	err = s.store.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, Error{Error: "Invalid ID"})
		return
	}
	c.JSON(http.StatusOK, "User deleted!")
}

// @Security ApiKeyAuth
// @Summary Get All Accounts
// @Description Fetch an account by account number
// @Tags account
// @Accept json
// @Produce json
// @Success 200 {object} types.Account
// @Failure 400 {object} Error "Bad Request"
// @Router /accounts [get]
func (s *Server) handleGetAllAccounts(c *gin.Context) {
	accounts, err := s.store.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

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
	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "go-bank-v2/docs"
	. "go-bank-v2/internal/types"
	"log"
	"net/http"
	"strconv"
)

func (s *Server) Run() {

	router := gin.Default()

	router.POST("/user", s.handleCreateUser)
	router.GET("/user/:id", withJWTAuth(s.handleGetUser, s.store))
	router.POST("/login", s.handleLogin)
	router.GET("/account/:id", withJWTAuth(s.handleGetAccountByID, s.store))
	router.POST("/account/:id", withJWTAuth(s.handleCreateAccount, s.store))
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

// @Security ApiKeyAuth
// @Summary Get Account by ID
// @Description Fetch an account by ID
// @Tags account
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} types.UserDto
// @Failure 400 {object} Error "Bad Request"
// @Router /account/{id} [get]
func (s *Server) handleGetAccountByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Error{Error: "Invalid ID"})
		return
	}

	if c.DefaultQuery("accountNumber", "") == "" {
		accounts, err := s.store.GetAccounts(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Error{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, accounts)
	} else {
		accNumStr := c.Query("accountNumber")
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
}

// @Security ApiKeyAuth
// @Summary Get User by ID
// @Description Fetch a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} types.Account
// @Failure 400 {object} Error "Bad Request"
// @Router /user/{id} [get]
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
	}
	c.JSON(http.StatusOK, userDto)
}

// @Security ApiKeyAuth
// @Summary Create an Account
// @Description Create an account for user with given ID
// @Tags account
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} types.Account
// @Failure 400 {object} Error "Bad Request"
// @Router /account/{id} [post]
func (s *Server) handleCreateAccount(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, Error{Error: "Method not allowed"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
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
// @Tags login
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

func (s *Server) handleCreateUser(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, Error{Error: "Error while creating a user"})
		return
	}
	firstName := c.Query("firstName")
	lastName := c.Query("lastName")
	password := c.Query("password")
	if firstName == "" || lastName == "" || password == "" {
		c.JSON(http.StatusMethodNotAllowed, Error{Error: "Error while creating user"})
		return
	}
	req := CreateUserRequest{
		FirstName: firstName,
		LastName:  lastName,
		Password:  password,
	}

	user, err := NewUser(req.FirstName, req.LastName, req.Password)
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

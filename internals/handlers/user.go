package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/aasif-10/Olx-API/internals/config"
	"github.com/aasif-10/Olx-API/internals/httpx"
	"github.com/aasif-10/Olx-API/internals/middlewares"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

var dummyHash = []byte("$2a$10$TZ79pRvNmPtFsJbSJSNaiucicFIXy7r.XgkwNux9P95hauXC9u4D2")

type user struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthHandler struct {
	db     *sql.DB
	logger *slog.Logger
	cfg    config.Config
}

func NewAuthHandler(db *sql.DB, logger *slog.Logger, cfg config.Config) *AuthHandler {
	return &AuthHandler{
		db:     db,
		logger: logger,
		cfg:    cfg,
	}
}

func (ah AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.RequestIDFromContext(ctx)

	log := ah.logger.With("request_id", requestId)

	var req SignupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error("Failed to decode", "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJSON)
		return
	}

	err = req.Validate()
	if err != nil {
		var verr *ValidationError
		errors.As(err, &verr)
		httpx.ValidationError(w, http.StatusBadRequest, err.Error(), httpx.CodeValidationFailed, verr.Field)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		log.Error("hashing failed", "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	row := ah.db.QueryRowContext(ctx,
		`INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING ID, CreatedAt`,
		req.Name, req.Email, hash)

	var u user
	err = row.Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			httpx.Error(w, http.StatusConflict, "email already taken", httpx.CodeConflict)
			return
		}

		log.Error("scanning failed", "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	out := SignupResponse{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
	}
	log.Info("new user registered", "user_id", out.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(out)
}

func (ah AuthHandler) Signin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.RequestIDFromContext(ctx)

	log := ah.logger.With("request_id", requestId)

	// 1. Store incoming request body
	var req SigninRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error("failed to decode", "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJSON)
		return
	}

	err = req.Validate()
	if err != nil {
		var verr *ValidationError
		errors.As(err, &verr)
		httpx.ValidationError(w, http.StatusBadRequest, err.Error(), httpx.CodeValidationFailed, verr.Field)
		return
	}

	// 2. DB Operation
	row := ah.db.QueryRowContext(ctx,
		`SELECT id, email, password FROM users where email = $1`, req.Email)

	var u user
	err = row.Scan(&u.ID, &u.Email, &u.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(req.Password))
			httpx.Error(w, http.StatusUnauthorized, "email or password don't match", httpx.CodeUnauthenticated)
			return
		}

		log.Error("find user by email failed", "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		log.Warn("password mismatch", "user_id", u.ID)
		httpx.Error(w, http.StatusUnauthorized, "email or password don't match", httpx.CodeUnauthenticated)
		return
	}

	tokenTTL := 24 * time.Hour
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   u.ID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(ah.cfg.JwtKey))
	if err != nil {
		log.Error("jwt sign failed", "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	out := SigninResponse{
		Token:     signed,
		ExpiresIn: int(tokenTTL.Seconds()),
	}
	log.Info("new user logged in", "user_id", u.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(out)

}

/*

type AuthHandler struct {
	db 		*sql.DB
	logger  *slog.Logger
}

func NewAuthHandler(db *sql.DB, logger *slog.Logger) *AuthHandler {
	return AuthHandler{
		db: db,
		logger: logger
	}
}

func (ah AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx:= req.Context()
	requestId:= middlewares.getRequestId(ctx)

	1. Store incoming json in a struct
	var req SignupRequest (dto)
	json.NewDecoder(req.Body).Decode(&req)

	2. Validation

	err:= req.Validate()
	if err!= nil {} //handle error

	3. Store in DB
	row = ah.db.QueryRowContext(ctx, `INSERT INTO users () VALUES ($1, ...) RETURNING ID, CreatedAt`,req.Name, ...)

	4. Store returining values from DB in "u"
	var u user
	_= row.Scan(&u.ID,&u.CreatedAt)

	5. Output json
	out:= SignupResponse{
		ID: u.ID,
		CreatedAt : u.CreatedAt,
	}

	w.Header().Set()..
	w.WriteHeader()..
	json.NewEncoder(w).Encode(out)
}

*/

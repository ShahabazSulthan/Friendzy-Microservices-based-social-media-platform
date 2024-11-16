package repository

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/requestmodels"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name    string
		input   *requestmodels.UserSignUpRequest
		stub    func(sqlmock.Sqlmock)
		wantErr error
	}{
		{
			name: "successfully inserted user",
			input: &requestmodels.UserSignUpRequest{
				Name:     "Shahabaz",
				UserName: "shahabaz123",
				Email:    "shahabaz@gmail.com",
				Password: "securepassword",
			},
			stub: func(s sqlmock.Sqlmock) {
				s.ExpectExec("INSERT INTO users").
					WithArgs("Shahabaz", "shahabaz123", "shahabaz@gmail.com", "securepassword").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: nil,
		},
		{
			name: "error inserting user",
			input: &requestmodels.UserSignUpRequest{
				Name:     "Shahabaz",
				UserName: "shahabaz123",
				Email:    "shahabaz@gmail.com",
				Password: "securepassword",
			},
			stub: func(s sqlmock.Sqlmock) {
				s.ExpectExec("INSERT INTO users").
					WithArgs("Shahabaz", "shahabaz123", "shahabaz@gmail.com", "securepassword").
					WillReturnError(fmt.Errorf("failed to insert user"))
			},
			wantErr: fmt.Errorf("failed to insert user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock DB and gorm DB
			mockDB, mockSql, _ := sqlmock.New()
			defer mockDB.Close()

			DB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})

			// Apply the stub
			tt.stub(mockSql)

			// Create the repository
			userRepository := NewUserRepo(DB)

			// Execute the method
			err := userRepository.CreateUser(tt.input)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			// Ensure all expectations were met
			err = mockSql.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestGetEmailAndUsernameByUserID(t *testing.T) {
	tests := []struct {
		name             string
		userID           int
		mockQuery        func(sqlmock.Sqlmock)
		expectedErr      error
		expectedEmail    string
		expectedUsername string
	}{
		{
			name:   "Successfully retrieved email and username",
			userID: 1,
			mockQuery: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"email", "user_name"}).
					AddRow("shahabazsultha4@gmaiil.com", "Shahabaz777")
				mock.ExpectQuery("SELECT email, user_name FROM users WHERE user_id = \\$1").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedErr:      nil,
			expectedEmail:    "shahabazsultha4@gmaiil.com",
			expectedUsername: "Shahabaz777",
		},
		{
			name:   "User not found",
			userID: 2,
			mockQuery: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"email", "user_name"})
				mock.ExpectQuery("SELECT email, user_name FROM users WHERE user_id = \\$1").
					WithArgs(2).
					WillReturnRows(rows)
			},
			expectedErr:      sql.ErrNoRows, // Expecting sql.ErrNoRows instead of gorm.ErrRecordNotFound
			expectedEmail:    "",
			expectedUsername: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock DB and gorm DB
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			assert.NoError(t, err)

			userRepo := &UserRepo{DB: gormDB}

			// Apply the stub
			tt.mockQuery(mock)

			// Execute the method
			email, username, err := userRepo.GetEmailAndUsernameByUserID(tt.userID)

			// Assertions
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedEmail, email)
				assert.Equal(t, tt.expectedUsername, username)
			}

			// Ensure all expectations were met
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestActivateUser(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		mockQuery   func(sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			name:  "Successfully activated user",
			email: "shahabazsultha4@gmaiil.com",
			mockQuery: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET status='active' WHERE email = \\$1").
					WithArgs("shahabazsultha4@gmaiil.com").
					WillReturnResult(sqlmock.NewResult(1, 1)) // Indicating one row was affected
			},
			expectedErr: nil,
		},
		{
			name:  "Error activating user",
			email: "shahabazsultha4@gmaiil.com",
			mockQuery: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET status='active' WHERE email = \\$1").
					WithArgs("shahabazsultha4@gmaiil.com").
					WillReturnError(fmt.Errorf("database error"))
			},
			expectedErr: fmt.Errorf("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock DB
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Setup GORM with the mocked DB
			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			assert.NoError(t, err)

			// Initialize UserRepo with the mocked DB
			userRepo := &UserRepo{DB: gormDB}

			// Apply the mock query
			tt.mockQuery(mock)

			// Execute the function
			err = userRepo.ActivateUser(tt.email)

			// Assertions
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			// Ensure all expectations were met
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestGetUserIDByEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		mockQuery   func(sqlmock.Sqlmock)
		expectedID  string
		expectedErr error
	}{
		{
			name:  "Successfully retrieved user ID by email",
			email: "shahabazsultha4@gmail.com",
			mockQuery: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow("12345")
				mock.ExpectQuery("SELECT id FROM users WHERE email=\\$1 And status=\\$2").
					WithArgs("shahabazsultha4@gmail.com", "active").
					WillReturnRows(rows)
			},
			expectedID:  "12345",
			expectedErr: nil,
		},
		{
			name:  "User not found",
			email: "shahabazsultha4@gmaiil.com",
			mockQuery: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id FROM users WHERE email=\\$1 And status=\\$2").
					WithArgs("shahabazsultha4@gmaiil.com", "active").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedID:  "",
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name:  "Database error",
			email: "shahabazsultha4@gmaiil.com",
			mockQuery: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id FROM users WHERE email=\\$1 And status=\\$2").
					WithArgs("shahabazsultha4@gmaiil.com", "active").
					WillReturnError(fmt.Errorf("database error"))
			},
			expectedID:  "",
			expectedErr: fmt.Errorf("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock DB
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Setup GORM with the mocked DB
			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			assert.NoError(t, err)

			// Initialize UserRepo with the mocked DB
			userRepo := &UserRepo{DB: gormDB}

			// Apply the mock query
			tt.mockQuery(mock)

			// Execute the function
			userID, err := userRepo.GetUserIDByEmail(tt.email)

			// Assertions
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Equal(t, tt.expectedID, userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, userID)
			}

			// Ensure all expectations were met
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestIsUserBlocked(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		mockQuery     func(sqlmock.Sqlmock)
		expectedBlocked bool
		expectedErr   error
	}{
		{
			name:   "User is blocked",
			userID: "1",
			mockQuery: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"status"}).AddRow("blocked")
				mock.ExpectQuery("SELECT status FROM users WHERE id=\\$1").
					WithArgs("1").
					WillReturnRows(rows)
			},
			expectedBlocked: true,
			expectedErr:     nil,
		},
		{
			name:   "User is not blocked",
			userID: "2",
			mockQuery: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"status"}).AddRow("active")
				mock.ExpectQuery("SELECT status FROM users WHERE id=\\$1").
					WithArgs("2").
					WillReturnRows(rows)
			},
			expectedBlocked: false,
			expectedErr:     nil,
		},
		{
			name:   "User not found",
			userID: "3",
			mockQuery: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT status FROM users WHERE id=\\$1").
					WithArgs("3").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedBlocked: false,
			expectedErr:     gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock DB
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Setup GORM with the mocked DB
			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			assert.NoError(t, err)

			// Initialize UserRepo with the mocked DB
			userRepo := &UserRepo{DB: gormDB}

			// Apply the mock query
			tt.mockQuery(mock)

			// Execute the function
			isBlocked, err := userRepo.IsUserBlocked(tt.userID)

			// Assertions
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.False(t, isBlocked)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBlocked, isBlocked)
			}

			// Ensure all expectations were met
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

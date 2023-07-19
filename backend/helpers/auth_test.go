package helpers

import (
	"reflect"
	"testing"

	"github.com/ryanozx/skillnet/models"
)

const (
	testUserID = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
)

func TestIsEmptyUserPass(t *testing.T) {
	type args struct {
		user *models.UserCredentials
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Non-empty username and password",
			args{
				user: &models.UserCredentials{
					Username: "Test user",
					Password: "12345",
				},
			},
			false,
		},
		{
			"Username with whitespace",
			args{
				user: &models.UserCredentials{
					Username: "    ",
					Password: "12345",
				},
			},
			true,
		},
		{
			"Password with whitespace",
			args{
				user: &models.UserCredentials{
					Username: "Test user",
					Password: "     ",
				},
			},
			true,
		},
		{
			"Empty username",
			args{
				user: &models.UserCredentials{
					Password: "12345",
				},
			},
			true,
		},
		{
			"Empty password",
			args{
				user: &models.UserCredentials{
					Username: "12345",
				},
			},
			true,
		},
		{
			"Empty username and password",
			args{
				user: &models.UserCredentials{},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmptyUserPass(tt.args.user); got != tt.want {
				t.Errorf("IsEmptyUserPass() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSignupUserCredsEmpty(t *testing.T) {
	type args struct {
		user *models.SignupUserCredentials
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Non-empty signup user creds",
			args{
				user: &models.SignupUserCredentials{
					UserCredentials: models.UserCredentials{
						Username: "TestUser",
						Password: "12345",
					},
					Email: "abc@def.com",
				},
			},
			false,
		},
		{
			"Email with whitespace",
			args{
				user: &models.SignupUserCredentials{
					UserCredentials: models.UserCredentials{
						Username: "TestUser",
						Password: "12345",
					},
					Email: "     ",
				},
			},
			true,
		},
		{
			"Empty user creds",
			args{
				user: &models.SignupUserCredentials{
					Email: "abc@def.com",
				},
			},
			true,
		},
		{
			"Empty email",
			args{
				user: &models.SignupUserCredentials{
					UserCredentials: models.UserCredentials{
						Username: "TestUser",
						Password: "12345",
					},
				},
			},
			true,
		},
		{
			"Empty user credentials and email",
			args{
				user: &models.SignupUserCredentials{},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSignupUserCredsEmpty(tt.args.user); got != tt.want {
				t.Errorf("IsSignupUserCredsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidSession(t *testing.T) {
	type args struct {
		session SessionGetter
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Valid session",
			args{
				session: &MockSessionStore{
					Values: map[interface{}]interface{}{UserIDKey: testUserID},
				},
			},
			true,
		},
		{
			"Invalid session",
			args{
				session: &MockSessionStore{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidSession(tt.args.session); got != tt.want {
				t.Errorf("IsValidSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockPostFormer struct {
	Params map[string](string)
}

func (pf *mockPostFormer) PostForm(key string) string {
	return pf.Params[key]
}

func TestExtractUserCredentials(t *testing.T) {
	type args struct {
		ctx postFormer
	}
	tests := []struct {
		name string
		args args
		want *models.UserCredentials
	}{
		{
			"Username and password present",
			args{
				ctx: &mockPostFormer{
					Params: map[string]string{"username": "TestUser", "password": "12345"},
				},
			},
			&models.UserCredentials{
				Username: "TestUser",
				Password: "12345",
			},
		},
		{
			"Username present",
			args{
				ctx: &mockPostFormer{
					Params: map[string]string{"username": "TestUser"},
				},
			},
			&models.UserCredentials{
				Username: "TestUser",
			},
		},
		{
			"Password present",
			args{
				ctx: &mockPostFormer{
					Params: map[string]string{"password": "12345"},
				},
			},
			&models.UserCredentials{
				Password: "12345",
			},
		},
		{
			"Username and password not present",
			args{
				ctx: &mockPostFormer{
					Params: map[string]string{},
				},
			},
			&models.UserCredentials{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractUserCredentials(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractUserCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractSignupUserCredentials(t *testing.T) {
	type args struct {
		ctx postFormer
	}
	tests := []struct {
		name string
		args args
		want *models.SignupUserCredentials
	}{
		{
			"User credentials and email present",
			args{
				ctx: &mockPostFormer{
					Params: map[string]string{"username": "TestUser", "password": "12345", "email": "abc@def.com"},
				},
			},
			&models.SignupUserCredentials{
				UserCredentials: models.UserCredentials{
					Username: "TestUser",
					Password: "12345",
				},
				Email: "abc@def.com",
			},
		},
		{
			"User credentials present",
			args{
				ctx: &mockPostFormer{
					Params: map[string]string{"username": "TestUser", "password": "12345"},
				},
			},
			&models.SignupUserCredentials{
				UserCredentials: models.UserCredentials{
					Username: "TestUser",
					Password: "12345",
				},
			},
		},
		{
			"Email present",
			args{
				ctx: &mockPostFormer{
					Params: map[string]string{"email": "abc@def.com"},
				},
			},
			&models.SignupUserCredentials{
				Email: "abc@def.com",
			},
		},
		{
			"User credentials and email not present",
			args{
				ctx: &mockPostFormer{
					Params: map[string]string{},
				},
			},
			&models.SignupUserCredentials{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractSignupUserCredentials(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractSignupUserCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

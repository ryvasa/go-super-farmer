package seeders

import (
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"golang.org/x/crypto/bcrypt"
)

func password() string {
	bytes, err := bcrypt.GenerateFromPassword([]byte("password"), 10)
	if err != nil {
		logrus.Log.Errorf("error generating password: %v", err)
	}
	password := string(bytes)
	return password
}

var users = []domain.User{
	{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: password(),
		ID:       uuid.New(),
		RoleID:   1,
	},
	{
		Name:     "Jane Doe",
		Email:    "jane.doe@example.com",
		Password: password(),
		ID:       uuid.New(),
		RoleID:   1,
	},
	{
		Name:     "Bob Smith",
		Email:    "bob.smith@example.com",
		Password: password(),
		ID:       uuid.New(),
		RoleID:   1,
	},
	{
		Name:     "Alice Johnson",
		Email:    "alice.johnson@example.com",
		Password: password(),
		ID:       uuid.New(),
		RoleID:   2,
	},
	{
		Name:     "Bob Martin",
		Email:    "bob.martin@example.com",
		Password: password(),
		ID:       uuid.New(),
		RoleID:   2,
	},
}

var commodities = []domain.Commodity{
	{
		Name:        "Wheat",
		ID:          uuid.New(),
		Description: "Wheat description",
		Code:        "WHT01",
		Duration:    time.Duration(3000 * time.Hour),
	},
	{
		Name:        "Corn",
		ID:          uuid.New(),
		Description: "Corn description",
		Code:        "CORN01",
		Duration:    time.Duration(3000 * time.Hour),
	},
	{
		Name:        "Rice",
		ID:          uuid.New(),
		Description: "Rice description",
		Code:        "RICE01",
		Duration:    time.Duration(3000 * time.Hour),
	},
	{
		Name:        "Soybean",
		ID:          uuid.New(),
		Description: "Soybean description",
		Code:        "SOY01",
		Duration:    time.Duration(3000 * time.Hour),
	},
}

var roles = []domain.Role{
	{ID: 1, Name: "Admin"},
	{ID: 2, Name: "Farmer"},
}

var provinces = []domain.Province{
	{Name: "Nanggroe Aceh Darussalam"},
	{Name: "Sumatera Utara"},
	{Name: "Sumatera Selatan"},
	{Name: "Sumatera Barat"},
	{Name: "Bengkulu"},
	{Name: "Riau"},
	{Name: "Kepulauan Riau"},
	{Name: "Jambi"},
	{Name: "Lampung"},
	{Name: "Bangka Belitung"},
	{Name: "Kalimantan Barat"},
	{Name: "Kalimantan Timur"},
	{Name: "Kalimantan Selatan"},
	{Name: "Kalimantan Tengah"},
	{Name: "Kalimantan Utara"},
	{Name: "Banten"},
	{Name: "DKI Jakarta"},
	{Name: "Jawa Barat"},
	{Name: "Jawa Tengah"},
	{Name: "Daerah Istimewa Yogyakarta"},
	{Name: "Jawa Timur"},
	{Name: "Bali"},
	{Name: "Nusa Tenggara Timur"},
	{Name: "Nusa Tenggara Barat"},
	{Name: "Gorontalo"},
	{Name: "Sulawesi Barat"},
	{Name: "Sulawesi Tengah"},
	{Name: "Sulawesi Utara"},
	{Name: "Sulawesi Tenggara"},
	{Name: "Sulawesi Selatan"},
	{Name: "Maluku Utara"},
	{Name: "Maluku"},
	{Name: "Papua Barat"},
	{Name: "Papua"},
	{Name: "Papua Tengah"},
	{Name: "Papua Pegunungan"},
	{Name: "Papua Selatan"},
	{Name: "Papua Barat Daya"},
}

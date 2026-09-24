package helper

import (
	"errors"
	"reflect"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		n := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if n == "" || n == "-" {
			return f.Name
		}
		return n
	})
	_ = v.RegisterValidation("nospace", func(fl validator.FieldLevel) bool { return !strings.ContainsAny(fl.Field().String(), " \t\n\r") })
	_ = v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		for _, r := range fl.Field().String() {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '.' && r != '_' {
				return false
			}
		}
		return true
	})
	_ = v.RegisterValidation("strongpassword", func(fl validator.FieldLevel) bool { return passwordStrength(fl.Field().String()) == "" })
	_ = v.RegisterValidation("nim", func(fl validator.FieldLevel) bool {
		for _, r := range fl.Field().String() {
			if r < '0' || r > '9' {
				return false
			}
		}
		return true
	})
	return v
}

func ValidateStruct(value any) map[string]string {
	err := validate.Struct(value)
	if err == nil {
		return nil
	}
	var invalid *validator.InvalidValidationError
	if errors.As(err, &invalid) {
		return map[string]string{"_": "objek yang divalidasi tidak sah"}
	}
	var fields validator.ValidationErrors
	if !errors.As(err, &fields) {
		return map[string]string{"_": "validasi gagal"}
	}
	result := map[string]string{}
	for _, field := range fields {
		if _, ok := result[field.Field()]; !ok {
			result[field.Field()] = messageFor(field)
		}
	}
	return result
}

func messageFor(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "wajib diisi"
	case "email":
		return "format email tidak valid"
	case "min":
		return "minimal " + fe.Param() + " karakter"
	case "max":
		return "maksimal " + fe.Param() + " karakter"
	case "alphanum":
		return "hanya boleh berisi huruf dan angka"
	case "nospace":
		return "tidak boleh mengandung spasi"
	case "username":
		return "hanya boleh huruf, angka, titik, dan garis bawah"
	case "strongpassword":
		return passwordStrength(fe.Value().(string))
	case "nim":
		return "NIM hanya boleh berisi angka"
	default:
		return "tidak memenuhi aturan " + fe.Tag()
	}
}

func passwordStrength(password string) string {
	if len(password) < 8 {
		return "minimal 8 karakter"
	}
	var letter, digit bool
	for _, r := range password {
		letter = letter || unicode.IsLetter(r)
		digit = digit || unicode.IsDigit(r)
	}
	if !letter || !digit {
		return "harus memuat huruf dan angka"
	}
	switch strings.ToLower(password) {
	case "password1", "12345678", "qwerty123", "admin123", "password123":
		return "password terlalu umum"
	}
	return ""
}

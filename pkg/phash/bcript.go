package phash

import "golang.org/x/crypto/bcrypt"

// GenerateFromPassword 生成密码哈希值
func GenerateFromPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	return string(b), err
}

// VerifyPassword 输入密码校验
func VerifyPassword(inputPassword, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}

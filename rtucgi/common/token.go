package common

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	secretKey = "yeda"
)

func GenerateToken() string {
	// 创建一个Token对象
	token := jwt.New(jwt.SigningMethodHS256)

	//jwt.SigningMethodHS256

	// 设置Token的自定义声明
	claims := token.Claims.(jwt.MapClaims)
	user, _ := GetUserInfo()
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 设置Token的过期时间

	// 使用密钥对Token进行签名，生成最终的Token字符串
	tokenString, _ := token.SignedString([]byte(secretKey))

	return tokenString
}

func IsTokenValid() bool {
	return ParseToken(GetAnyQuery("token"))
}

func ParseToken(tokenString string) bool {
	// 解析Token字符串
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return false
	}

	// 验证Token的签名方法是否有效
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return false
	}

	// 返回Token中的声明部分
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//claims
		user, _ := GetUserInfo()
		if claims["username"] == user.Username {
			return true
		}
	}
	//token.Claims

	return false
}

/*

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	// 解析Token字符串
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	//token.Valid

	if err != nil {
		return nil, err
	}

	// 验证Token的签名方法是否有效
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("无效的签名方法：%v", token.Header["alg"])
	}

	// 返回Token中的声明部分
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//claims
		a := claims["username"]
		return claims, nil
	}
	//token.Claims

	return nil, fmt.Errorf("无效的Token")
}
*/

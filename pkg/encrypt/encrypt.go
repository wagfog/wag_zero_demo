package encrypt

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/zeromicro/go-zero/core/codec"
)

const (
	passwordEncryptSeed = "(swag)@#$" //加密密码或其他敏感信息的种子值
	//AES（Advanced Encryption Standard）是一种对称密钥加密算法，用于加密和解密数据。它是一种高度安全的加密算法，被广泛用于保护敏感信息，包括银行交易、网络通信和数据存储。
	mobileAesKey = "5A2E746B08D846502F37A6E2D85D583B" //移动端应用程序中执行AES加密算法所需的密钥。
)

func EncPassword(password string) string {
	return Md5Sum([]byte(strings.TrimSpace(password + passwordEncryptSeed)))
}

func EncMoblie(mobile string) (string, error) {
	//用给定的key进行加密
	data, err := codec.EcbEncrypt([]byte(mobileAesKey), []byte(mobile))
	if err != nil {
		return "", err
	}
	//将数据进行Base64编码
	//Base64编码是一种用64个字符来表示任意二进制数据的方法，通常用于在数据传输过程中将二进制数据表示为文本数据，或者在需要时进行数据混淆。
	// /这种编码方法常常用在数据传输中，特别是在网络传输过程中
	return base64.StdEncoding.EncodeToString(data), nil
}

func DecMobile(mobile string) (string, error) {
	originalData, err := base64.StdEncoding.DecodeString(mobile)
	if err != nil {
		return "", err
	}
	data, err := codec.EcbDecrypt([]byte(mobileAesKey), originalData)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Md5Sum(data []byte) string {

	return hex.EncodeToString(byte16ToBytes(md5.Sum(data)))
}

// 这个函数的作用似乎是将长度为16的固定大小的字节数组转换为一个切片。该函数接受一个长度为16的字节数组作为输入，并将其转换为一个长度为16的切片
// 因为MD5返回的是长为16的数组
func byte16ToBytes(in [16]byte) []byte {
	tmp := make([]byte, 16)
	for _, value := range in {
		tmp = append(tmp, value)
	}
	return tmp[16:]
}

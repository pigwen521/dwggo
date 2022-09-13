package str

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
func Sha256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
func Sha1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// DES加密 des3
func Encrypt_DES3(origData, key []byte, iv []byte) ([]byte, error) {
	key = genkey_DES3(key)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	origData = pKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return Base64_encode_byte(crypted), nil
}

// DES解密 des3
func Decrypt_DES3(crypted, key, iv []byte) ([]byte, error) {
	key = genkey_DES3((key))
	crypted_byte := Base64_decode_byte(crypted)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted_byte))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted_byte)
	origData = pKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}
func genkey_DES3(key []byte) []byte {
	if len(key) < 24 {
		repeat_byte := bytes.Repeat([]byte("0"), 24)
		key = append(key, repeat_byte...)
	}
	return key[:24]
}

// AES加密 CBS CS5
func Encrypt_AES_CBC_CS5(origData_byte, key_byte []byte, iv []byte) ([]byte, error) {
	aesBlockEncrypter, err := aes.NewCipher(key_byte)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	content := pKCS5Padding(origData_byte, aesBlockEncrypter.BlockSize())
	encrypted := make([]byte, len(content))

	aesEncrypter := cipher.NewCBCEncrypter(aesBlockEncrypter, iv) //IV
	aesEncrypter.CryptBlocks(encrypted, content)
	return Base64_encode_byte(encrypted), nil
}

// AES解密 CBS CS5
func Decrypt_AES_CBC_CS5(encryptData, key_byte []byte, iv []byte) (data []byte, err error) {
	encryptData = Base64_decode_byte(encryptData)

	decrypted := make([]byte, len(encryptData))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher(key_byte)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	aesDecrypter := cipher.NewCBCDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.CryptBlocks(decrypted, encryptData)
	return pKCS5UnPadding(decrypted), nil
}

func pKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}
func pKCS5UnPadding(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

// 密码加密
func PasswordHash(passwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return string(bytes), err
}

// 密码解密
func PasswordVerify(passwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))
	return err == nil
}

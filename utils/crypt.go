package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

type aesCrypt func(origData []byte, key []byte) (crypted []byte, err error)

type AesEncryptMode int

const (
	EN_CBC AesEncryptMode = iota
	EN_ECB AesEncryptMode 
	EN_CFB AesEncryptMode
)

func (AesEncryptMode) Int() int {
	return AesEncryptMode.(int)
}

type AesDecryptMode int

const (
	DE_CBC AesDecryptMode = iota
	DE_ECB AesDecryptMode 
	DE_CFB AesDecryptMode
)

func (AesDecryptMode) Int() int {
	return AesDecryptMode.(int)
}

// =================== CBC ======================
func AesEncryptCBC(origData []byte, key []byte) (encrypted []byte, err error) {
	// 分組秘鑰
	// NewCipher該函數限制了輸入k的長度必須為16, 24或者32
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()                              // 獲取秘鑰塊的長度
	origData = pkcs5Padding(origData, blockSize)                // 補全碼
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted = make([]byte, len(origData))                     // 創建數組
	blockMode.CryptBlocks(encrypted, origData)                  // 加密
	return encrypted, nil
}
func AesDecryptCBC(encrypted []byte, key []byte) (decrypted []byte, err error) {
	block, err := aes.NewCipher(key) // 分組秘鑰
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()                              // 獲取秘鑰塊的長度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted = make([]byte, len(encrypted))                    // 創建數組
	blockMode.CryptBlocks(decrypted, encrypted)                 // 解密
	decrypted = pkcs5UnPadding(decrypted)                       // 去除補全碼
	return decrypted, nil
}
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// =================== ECB ======================
func AesEncryptECB(origData []byte, key []byte) (encrypted []byte, err error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分組分塊加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted, nil
}
func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte, err error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypted = make([]byte, len(encrypted))
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	if trim <= 0 {
		return []byte{}, errors.New("AesDecryptECB Fail")
	}
	return decrypted[:trim], nil
}

// =================== CFB ======================
func AesEncryptCFB(origData []byte, key []byte) (encrypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted, nil
}
func AesDecryptCFB(encrypted []byte, key []byte) (decrypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(encrypted) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted, nil
}

//=================== MD5 ======================
func Md5Encrypt(origData ...string) string {
	str := ""
	for _, v := range origData {
		str += v
	}
	md5byte := md5.Sum([]byte(str))
	return hex.EncodeToString(md5byte[:])
}

/* 使用base64封裝其餘AES加密方式 */
func AesBase64Encrypt(origData []byte, key []byte,mode AesEncryptMode) (got string, err error) {
	// gotEncrypted, err := aesFunc(origData, key)
	var gotEncrypted []byte
	switch mode {
	case EN_CBC:
		gotEncrypted , err = AesEncryptCBC(b,key)
	case EN_ECB:
		gotEncrypted , err = AesEncryptECB(b,key)
	case EN_CFB:
		gotEncrypted , err = AesEncryptCFB(b,key)
	default:
		return nil, errors.New("mode error")
	}
	if err != nil {
		return "", err
	}
	got = base64.StdEncoding.EncodeToString(gotEncrypted)
	return got, err
}

/* 使用base64封裝其餘AES解密方式 */
func AesBase64Decrypt(origData string, key []byte,mode AesDecryptMode) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(origData)
	if err != nil {
		return nil, err
	}
	switch mode {
	case DE_CBC:
		return AesDecryptCBC(b,key)
	case DE_ECB:
		return AesDecryptECB(b,key)
	case DE_CFB:
		return AesDecryptCFB(b,key)
	default:
		return nil, errors.New("mode error")
	}
}

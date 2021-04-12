package util

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"time"
)

const AESKey = "idfodfdfpofofgfx"

func AesEncrypt(origData, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()

	origData = PKCS5Padding(origData, blockSize)

	// origData = ZeroPadding(origData, block.BlockSize())

	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])

	crypted := make([]byte, len(origData))

	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以

	// crypted := origData

	blockMode.CryptBlocks(crypted, origData)

	return crypted, nil

}

func AesDecrypt(crypted, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)

	if err != nil {

		return nil, err

	}

	blockSize := block.BlockSize()

	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])

	origData := make([]byte, len(crypted))

	// origData := crypted

	blockMode.CryptBlocks(origData, crypted)

	origData = PKCS5UnPadding(origData)

	// origData = ZeroUnPadding(origData)

	return origData, nil

}
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize //需要padding的数目
	//只要少于256就能放到一个byte中，默认的blockSize=16(即采用16*8=128, AES-128长的密钥)
	//最少填充1个byte，如果原文刚好是blocksize的整数倍，则再填充一个blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding) //生成填充的文本
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
func MD5(str string) (string, error) {
	md := crypto.MD5.New()
	if _, e := md.Write([]byte(str)); e != nil {
		return "", e
	}
	return hex.EncodeToString(md.Sum(nil)), nil

}

func GetToken(id string) (string, error) {

	m := map[string]interface{}{
		"id":        id,
		"timestamp": time.Now().Unix(),
	}
	if b, e := json.Marshal(m); e == nil {
		if b, e = AesEncrypt(b, []byte(AESKey)); e == nil {
			return hex.EncodeToString(b), nil
		}
		return "", e

	} else {
		return "", e
	}

}
func Token(t string) (string, int64, error) {
	b, e := hex.DecodeString(t)
	if e != nil {
		return "", 0, e

	}

	b, e = AesDecrypt(b, []byte(AESKey))
	if e != nil {
		return "", 0, e
	}
	var m struct {
		Id        string `json:"id"`
		Timestamp int64  `json:"timestamp"`
	}
	if e := json.Unmarshal(b, &m); e != nil {
		return "", 0, e
	}

	return m.Id, m.Timestamp, nil

}

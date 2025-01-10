package main

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	mr "math/rand"
	"strconv"
	"time"
)

// LicenseInfo 表示license的内容结构
type LicenseInfo struct {
	AppID        string `json:"appId"`
	IssuedTime   int64  `json:"issuedTime"`
	NotBefore    int64  `json:"notBefore"`
	NotAfter     int64  `json:"notAfter"`
	CustomerInfo string `json:"customerInfo"`
}

func generateRsaKey(bits int) (privKey, pubKey string, err error) {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	// 导出私钥为PEM格式
	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})

	// 导出公钥为PEM格式
	pubBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	pubPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubBytes,
	})

	return string(privPem), string(pubPem), nil
}

// 解析RSA私钥
func parsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// 解析RSA公钥
func parsePublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

// PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// PKCS7 去除填充
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("invalid padding data")
	}
	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, fmt.Errorf("invalid padding size")
	}
	return data[:(length - unpadding)], nil
}

// GenerateLicense 生成license
func GenerateLicense(info *LicenseInfo, rsaPrivateKey string) (string, error) {
	// 1. 生成16位随机字符串作为AES密钥
	aesKey := MyUtil_getRandomString(16)

	// 2. 将license信息转换为JSON
	jsonData, err := json.Marshal(info)
	if err != nil {
		return "", fmt.Errorf("marshal license info: %v", err)
	}

	// 3. 使用AES加密JSON数据
	encData, err := AesUtil_encrypt(aesKey, string(jsonData))
	if err != nil {
		return "", fmt.Errorf("AES encrypt: %v", err)
	}

	// 4. 计算加密数据长度（16进制，2字符）
	encDataLenHex := fmt.Sprintf("%02x", len(encData))

	// 5. 使用RSA私钥对加密数据进行签名
	sign, err := RSAUtil_sign(rsaPrivateKey, encData)
	if err != nil {
		return "", fmt.Errorf("create signature: %v", err)
	}

	// 6. 组合最终的license字符串 (aesKey + encDataLength + encData + sign)
	license := fmt.Sprintf("%s%s%s%s", aesKey, encDataLenHex, encData, sign)
	return license, nil
}

// MyUtil_getRandomString 生成指定长度的随机字符串
func MyUtil_getRandomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[mr.Intn(len(chars))]
	}
	return string(result)
}

// AesUtil_encrypt AES加密
func AesUtil_encrypt(key, content string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 使用PKCS7填充
	blockSize := block.BlockSize()
	content = string(pkcs7Padding([]byte(content), blockSize))

	crypted := make([]byte, len(content))
	// ECB加密
	for bs, be := 0, blockSize; bs < len(content); bs, be = bs+blockSize, be+blockSize {
		block.Encrypt(crypted[bs:be], []byte(content)[bs:be])
	}

	// Base64编码
	return base64.StdEncoding.EncodeToString(crypted), nil
}

// AesUtil_decrypt AES解密
func AesUtil_decrypt(key, content string) (string, error) {
	// Base64解码
	crypted, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	decrypted := make([]byte, len(crypted))
	size := block.BlockSize()

	// ECB解密
	for bs, be := 0, size; bs < len(crypted); bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], crypted[bs:be])
	}

	// 去除PKCS7填充
	decrypted, err = pkcs7UnPadding(decrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// RSAUtil_sign RSA签名
func RSAUtil_sign(privateKeyPEM string, data string) (string, error) {
	privateKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return "", err
	}

	hashed := sha1.Sum([]byte(data))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hashed[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// RSAUtil_verifySign RSA验签
func RSAUtil_verifySign(publicKeyPEM string, data string, sign string) (bool, error) {
	publicKey, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return false, err
	}

	signature, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}

	hashed := sha1.Sum([]byte(data))
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hashed[:], signature)
	return err == nil, err
}

// VerifyLicense 验证license
func VerifyLicense(license string, rsaPublicKey string, expectedAppID string) (*LicenseInfo, error) {
	if len(license) < 18 {
		return nil, fmt.Errorf("invalid license format")
	}

	// 1. 解析license各部分
	aesKey := license[:16]          // 16字节的AES密钥
	encDataLenHex := license[16:18] // 2字节的长度信息

	// 解析加密数据长度
	dataLen, err := strconv.ParseInt(encDataLenHex, 16, 64)
	if err != nil {
		return nil, fmt.Errorf("decode data length: %v", err)
	}

	// 提取加密数据和签名
	encDataStart := 18
	encDataEnd := encDataStart + int(dataLen)
	if encDataEnd >= len(license) {
		return nil, fmt.Errorf("invalid license format")
	}

	encData := license[encDataStart:encDataEnd]
	sign := license[encDataEnd:]

	// 2. 验证RSA签名
	valid, err := RSAUtil_verifySign(rsaPublicKey, encData, sign)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid signature")
	}

	// 3. 解密数据
	jsonData, err := AesUtil_decrypt(aesKey, encData)
	if err != nil {
		return nil, fmt.Errorf("AES decrypt: %v", err)
	}

	// 4. 解析JSON数据
	var info LicenseInfo
	if err := json.Unmarshal([]byte(jsonData), &info); err != nil {
		return nil, fmt.Errorf("unmarshal license info: %v", err)
	}

	// 5. 验证appId和时间
	if info.AppID != expectedAppID {
		return nil, fmt.Errorf("invalid app ID")
	}

	now := time.Now().UnixMilli()
	if now < info.NotBefore || now > info.NotAfter {
		return nil, fmt.Errorf("license expired or not yet valid")
	}

	return &info, nil
}

func main() {
	privateKey, publicKey, err := generateRsaKey(2048)
	if err != nil {
		panic(err)
	}

	// 创建license信息
	info := &LicenseInfo{
		AppID:        "com.example.app",
		IssuedTime:   time.Now().UnixMilli(),
		NotBefore:    time.Now().UnixMilli(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour).UnixMilli(),
		CustomerInfo: "Customer Name",
	}

	// 生成license
	license, err := GenerateLicense(info, privateKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("Generated License:", license)
	fmt.Println("License length:", len(license))

	// 分段打印license内容
	fmt.Println("AES Key:", license[:16])
	fmt.Println("Length Hex:", license[16:18])
	dataLen, _ := strconv.ParseInt(license[16:18], 16, 64)
	fmt.Println("Encrypted Data:", license[18:18+dataLen])
	fmt.Println("Signature:", license[18+dataLen:])

	// 验证license
	verifiedInfo, err := VerifyLicense(license, publicKey, "com.example.app")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Verified License Info: %+v\n", verifiedInfo)
}

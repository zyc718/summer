package example

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/pkg/errors"
)

//私钥
var privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDLPVEhoZSXNJ45U0mEMndBlJ68UY6KUiOVrR5jPPp+Cx6Eq1+F
eWPLmC/T5G37b0GeKG+tYJb5eBY4M99QGKBK8nj0Zcu4tTEAlQPwUbkbGafwPdCY
e1NY5QcwpLQzabYes/TOESdLZHi6A/WGmmKMdHc+sn1333elNpxrJUB7FQIDAQAB
AoGAXy7vYUXQVmRhOc3E33HXIlKdaOr1S9ieK8oxMh7r3b4NY+ryyIsKbt5uf9k6
nQgE/jJH4zYaXumb1mSM0HFIGBtw7aoaEl7SGYe9thJGmwKZCdc95MjHyGmBIRi2
fG4T/qUNe0DbZwFGXnIm8F3yF6g8F78eIOD/svMqdqPAC4ECQQD1SNpcR2R1dOeE
DBXoSmmpkpAyCwdzHwd/khSvOa/WPyXRByms2ZZOI6AGw7TkJRQxiB+RMWM2eyRY
2b4SBXrhAkEA1B5Adl1h8fqqqYiZPxCEK91jlIChd0OQJin61uWXGMlln49IG54W
ZoQ2QE8IPn1Hd0HT6YYyz3+rSboWOV3atQJBANu/izHFHDFGrOvWUAIuOH+dOOY8
j04J7JPT8ggSLIBLTrv4KNQck9YpgILO7s6+kVrW00Em9/WlWSjo2qoWksECQQDB
+hhBJgyX2P+QodZikZwM8RxLhYYjJqn//IvjUXnntOU2ETWD7AHYJjfmf1+upaph
KNW9zHdSwhHGDmKce3OxAkAkYBKj+JRe0nnH4qT9v18oT0JQC3FxCvwB3LGWFoDo
mOjU0g4Ctq21Gg/H48GhgUWxqOMfteImIz+g5W+Z2zJq
-----END RSA PRIVATE KEY-----
`)

//公钥
var publicKey = []byte(`-----BEGIN RSA PUBLIC KEY-----
MIGJAoGBAMs9USGhlJc0njlTSYQyd0GUnrxRjopSI5WtHmM8+n4LHoSrX4V5Y8uY
L9PkbftvQZ4ob61glvl4Fjgz31AYoEryePRly7i1MQCVA/BRuRsZp/A90Jh7U1jl
BzCktDNpth6z9M4RJ0tkeLoD9YaaYox0dz6yfXffd6U2nGslQHsVAgMBAAE=
-----END RSA PUBLIC KEY-----
`)

// 加密
func RsaEncrypt(origData []byte) (string, error) {
	block, _ := pem.Decode(publicKey) //将密钥解析成公钥实例
	if block == nil {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKCS1PublicKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	if err != nil {
		return "", err
	}
	//pub := pubInterface.(*rsa.PublicKey)

	encodeByte, errEncrypt := rsa.EncryptPKCS1v15(rand.Reader, pubInterface, origData) //RSA算法加密
	if errEncrypt != nil {
		err = errors.Wrap(errEncrypt, "Eecoding DecryptPKCS1v15")
		return "", errEncrypt
	}
	return base64.StdEncoding.EncodeToString(encodeByte), nil
}

// 解密
func RsaDecrypt(ciphertext string) (string, error) {
	block, _ := pem.Decode(privateKey) //将密钥解析成私钥实例
	if block == nil {
		return "", errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	if err != nil {
		return "", err
	}

	ciphertextBase, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	encodeByte, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertextBase) //RSA算法解密
	if err != nil {
		return "", err
	}
	return string(encodeByte), nil

}

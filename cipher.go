package main

type Cipher interface {
	//	Encrypt(strMsg string) ([]byte, error)
	Decrypt(src []byte) (decrypted []byte, err error)
}

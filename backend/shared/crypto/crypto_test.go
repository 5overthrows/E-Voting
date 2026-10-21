package crypto

import (
	"testing"
)

func TestPasswordHashing(t *testing.T) {
	pass := "mySecretPass123"
	hash, err := HashPassword(pass)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if !CheckPasswordHash(pass, hash) {
		t.Errorf("Password hash verification failed")
	}

	if CheckPasswordHash("wrongpass", hash) {
		t.Errorf("Password hash matched incorrect password")
	}
}

func TestRSAEncryptionDecryption(t *testing.T) {
	pub, priv, err := GenerateRSAKeypair()
	if err != nil {
		t.Fatalf("Failed to generate RSA keypair: %v", err)
	}

	plaintext := []byte("candidate_abc_123")
	ciphertext, err := EncryptRSA(pub, plaintext)
	if err != nil {
		t.Fatalf("RSA Encryption failed: %v", err)
	}

	decrypted, err := DecryptRSA(priv, ciphertext)
	if err != nil {
		t.Fatalf("RSA Decryption failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("Decrypted text '%s' does not match plaintext '%s'", string(decrypted), string(plaintext))
	}
}

func TestAESEncryptionDecryption(t *testing.T) {
	passphrase := "counting_authority_key"
	salt, err := GenerateRandomSalt()
	if err != nil {
		t.Fatalf("Failed to generate salt: %v", err)
	}

	key := DeriveKey(passphrase, salt)
	data := []byte("-----BEGIN RSA PRIVATE KEY-----\nMOCK PRIVATE KEY DATA\n-----END RSA PRIVATE KEY-----")

	ciphertext, err := EncryptAES(key, data)
	if err != nil {
		t.Fatalf("AES Encryption failed: %v", err)
	}

	decrypted, err := DecryptAES(key, ciphertext)
	if err != nil {
		t.Fatalf("AES Decryption failed: %v", err)
	}

	if string(decrypted) != string(data) {
		t.Errorf("Decrypted data does not match original data")
	}

	// Verify decryption fails with wrong passphrase
	wrongKey := DeriveKey("wrong_passphrase", salt)
	_, err = DecryptAES(wrongKey, ciphertext)
	if err == nil {
		t.Errorf("AES Decryption succeeded with incorrect key")
	}
}

func TestRSASigningVerification(t *testing.T) {
	pub, priv, err := GenerateRSAKeypair()
	if err != nil {
		t.Fatalf("Failed to generate keypair: %v", err)
	}

	msg := []byte("merkle_root_hash_value_12345")
	sig, err := SignRSA(priv, msg)
	if err != nil {
		t.Fatalf("RSA signing failed: %v", err)
	}

	err = VerifyRSASignature(pub, msg, sig)
	if err != nil {
		t.Errorf("RSA signature verification failed: %v", err)
	}

	// Verify fails with modified message
	err = VerifyRSASignature(pub, []byte("modified_msg"), sig)
	if err == nil {
		t.Errorf("Signature verified successfully for modified message")
	}
}

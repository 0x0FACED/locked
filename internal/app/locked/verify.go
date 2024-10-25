package locked

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/term"
)

func initApp() error {
	fmt.Print("Enter master password: ")
	var masterPass string
	fmt.Scanln(&masterPass)

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	hash := pbkdf2.Key([]byte(masterPass), salt, 10000, 32, sha256.New)

	if err := os.MkdirAll("secrets", os.ModePerm); err != nil {
		return err
	}

	filePath := filepath.Join("secrets", "master_hash")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(append(salt, hash...)); err != nil {
		return err
	}

	if err := os.Chmod(filePath, 0600); err != nil {
		return err
	}

	fmt.Println("Master password set up successfully.")
	return nil
}

func readPasswordData() ([]byte, []byte, error) {
	data, err := os.ReadFile(filepath.Join("secrets", "master_hash"))
	if err != nil {
		return []byte{}, []byte{}, err
	}

	salt := data[:16]
	hash := data[16:]

	return salt, hash, nil
}

func verify() {
	salt, storedHash, err := readPasswordData()
	if err != nil {
		fmt.Println("Error reading password data:", err)
		return
	}

	for {
		inputPassword, err := requestPassword()
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		inputHash := pbkdf2.Key([]byte(inputPassword), salt, 10000, 32, sha256.New)

		if !compareHashes(storedHash, inputHash) {
			fmt.Println("Incorrect password, try again")
		} else {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func compareHashes(hash1, hash2 []byte) bool {
	if len(hash1) != len(hash2) {
		return false
	}
	for i := range hash1 {
		if hash1[i] != hash2[i] {
			return false
		}
	}
	return true
}

// запрашиваем ввод пароля
func requestPassword() (string, error) {
	fmt.Print("Enter master password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println() // Перенос строки после ввода
	return string(bytePassword), nil
}

package kvp

import (
	"encoding/json"

	"os"
	"strings"

	"fmt"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"bytes"
	"io/ioutil"
	"crypto/aes"
	"crypto/cipher"
)

type envvars map[string]string

var s3Svc *s3.S3
var cfg aws.Config

var aesKey [32]byte

func init(){
	copy(aesKey[:], []byte("ahnahb8yg6f6rybnhuvvvvtgtg88hF56"))
}

func EnvPairs(secretsNames []string) ([]string, error) {
	initialVars := envVarsFromOS(os.Environ())

	awsVars, err := envVarsFromAWS(secretsNames)
	if err != nil {
		return nil, err
	}

	for k, v := range awsVars {
		initialVars[k] = v
	}

	var out []string
	for k, v := range initialVars {
		out = append(out, fmt.Sprintf("%s=%v", k, v))
	}

	return out, nil
}

func envVarsFromOS(pairs []string) envvars {

	e := make(map[string]string)

	for _, v := range pairs {
		kvpair := strings.Split(v, "=")
		if len(kvpair) == 2 {
			e[kvpair[0]] = kvpair[1]
		}
	}

	return envvars(e)
}

func envVarsFromAWS(secretsNames []string) (envvars, error) {

	var err error
	cfg, err = external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	if region := os.Getenv("SEEKRITS_REGION"); region != ""{
		cfg.Region = region
	}

	smSvc := secretsmanager.New(cfg)

	outData := make(map[string]string)

	for _, secretName := range secretsNames {

		req := smSvc.GetSecretValueRequest(&secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretName),
		})

		resp, err := req.Send()
		if err != nil {
			return nil, err
		}

		var respData map[string]string

		err = json.Unmarshal([]byte(*resp.SecretString), &respData)
		if err != nil {
			return nil, err
		}

		for k, v := range respData {
			if strings.HasPrefix(k, "s3://"){
				var err error
				k, v, err = decryptKeyFile(k, v)
				if err != nil {
					return nil, err
				}
			}
			outData[k] = v
		}
	}

	return envvars(outData), nil
}

func decryptKeyFile(k, v string)(string, string, error){
	rawFileData, objectName, err := getEncryptedFile([]byte(k))
	if err != nil {
		return "", "", err
	}
	if len(v) != 32{
		return "", "", errors.New("aes encryption key mush be 32 bytes in length")
	}
	var key [32]byte
	copy(key[:], []byte(v))
	plain, err := decrypt(*rawFileData, &key)
	if err != nil {
		return "", "", err
	}

	return string(objectName), string(plain), nil
}

func getEncryptedFile(bucketPath []byte) (*[]byte, []byte,  error){
	bucketPath = bucketPath[5:]
	pathParts := bytes.Split(bucketPath, []byte{'/'})

	bucketName := pathParts[0]
	objectPath := bucketPath[len(bucketName):]
	objectName := pathParts[len(pathParts)-1]

	req := s3Svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(string(bucketName)),
		Key: aws.String(string(objectPath)),
	})

	resp, err := req.Send()
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return &output, objectName, nil
}

func decrypt(ciphertext []byte, key *[32]byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

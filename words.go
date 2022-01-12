package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"unicode"

	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const BUCKET = "dk.lundogbendsen.aws.buckets.spellingbee"
const KEY = "corncob_lowercase.txt"
const LOCALFILENAME = "./corncob_lowercase.txt"
const AWSREGION = "eu-west-1"

func exitErrorf(msg string, args ...interface{}) {
	log.Printf(msg+"\n", args...)
	os.Exit(1)
}

func GetMatchingWords(letters string, mandatoryChar rune) ([]string, error) {
	letters = strings.ToLower(letters)
	mandatoryChar = unicode.ToLower(mandatoryChar)

	err := verify(letters, mandatoryChar)
	if err != nil {
		return nil, err
	}

	dictionary := make([]string, 0)

	//if we're running on aws
	if isRunningOnAwsLambda() {
		session, svc := getSessionAndS3()
		ensureBucketAndDictionaryExists(session, svc)
		dictionary, err = getWordlistFromS3(session, svc)
		if err != nil {
			exitErrorf("could not read dictionary from bucket: %s", err.Error())
		}
	} else { //running locally
		dictionary, err = GetWordlistFromLocalFile()
		if err != nil {
			exitErrorf("could not read dictionary from local file: %s", err.Error())
		}
	}
	return getWordFromLetters(letters, mandatoryChar, dictionary), nil
}

//get all possible words from dictionary containing only letters and mandatoryChar
func getWordFromLetters(letters string, mandatoryChar rune, wordlist []string) []string {
	//make sure chars are unique and sorted
	letters = uniqueLettersSorted(letters)

	result := make([]string, 0)

	for i := 0; i < len(wordlist); i++ {
		word := uniqueLettersSorted(wordlist[i])

		match := true
		//se if each char in word is found in letters
		for _, char := range word {
			if !strings.ContainsRune(letters, char) {
				match = false
				break
			}
		}

		//does the word contain the mandatory char?
		if match {
			if strings.ContainsRune(wordlist[i], mandatoryChar) {
				result = append(result, wordlist[i])
			}
		}

	}
	return result
}

//reads words from local file and returns a string array with one line per slot
func GetWordlistFromLocalFile() ([]string, error) {
	fn := LOCALFILENAME
	dat, err := os.Open(fn) //"./corncob_lowercase.txt")

	if err != nil {
		return nil, errors.New("could not open dictionary. cause:" + err.Error())
	}

	wl := make([]string, 0)
	scanner := bufio.NewScanner(dat)
	for scanner.Scan() {
		wl = append(wl, scanner.Text())
	}
	dat.Close()
	return wl, nil
}

//reads words from file in s3 bucket and returns a string array with one line per slot
func getWordlistFromS3(session *session.Session, svc *s3.S3) ([]string, error) {
	requestInput := &s3.GetObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(KEY),
	}

	result, err := svc.GetObject(requestInput)

	if err != nil {
		return nil, errors.New("Unable get object in bucket.. cause:" + err.Error())
	}

	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	wordsString := fmt.Sprintf("%s", body)

	wl := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(wordsString))

	for scanner.Scan() {
		wl = append(wl, scanner.Text())
	}
	return wl, nil
}

func isRunningOnAwsLambda() bool {
	return os.Getenv("LAMBDA_TASK_ROOT") != ""
}

func getSessionAndS3() (*session.Session, *s3.S3) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AWSREGION)},
	)

	if err != nil {
		exitErrorf("error creating session", err)
	}

	// Create S3 service client
	return sess, s3.New(sess)

}

//if file cannot be retrieved from bucket a bucket is created and the dictionary file uploaded
func ensureBucketAndDictionaryExists(sess *session.Session, svc *s3.S3) {

	root := os.Getenv("LAMBDA_TASK_ROOT")

	_, err := svc.HeadObject(&s3.HeadObjectInput{Bucket: aws.String(BUCKET), Key: aws.String(KEY)})

	//we cannot retrieve the object (wordlist)
	if err != nil {
		//does bucket exist?
		_, err = svc.HeadBucket(&s3.HeadBucketInput{Bucket: aws.String(BUCKET)})

		//if not, create it
		if err != nil {
			_, err = svc.CreateBucket(&s3.CreateBucketInput{
				Bucket: aws.String(BUCKET),
			})
			log.Print("bucket created")

			err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
				Bucket: aws.String(BUCKET),
			})
			log.Print("bucket created...done waiting")
		}

		//upload file
		fn := root + "/corncob_lowercase.txt"
		file, err2 := os.Open(fn)

		if err2 != nil {
			exitErrorf("error open words file", err2)
		}

		defer file.Close()

		uploader := s3manager.NewUploader(sess)

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(BUCKET),
			Key:    aws.String(KEY),
			Body:   file,
		})

		if err != nil {
			// Print the error and exit.
			exitErrorf("Unable to upload %q to %q, %v", "corncob_lowercase.txt", BUCKET, err)
		}

		log.Printf("Successfully uploaded %q to %q\n", "corncob_lowercase.txt", BUCKET)

	}
}

/* verifies preconditions */
func verify(letters string, mandatoryChar rune) error {
	uniqueLettersSorted := uniqueLettersSorted(letters)

	if len(uniqueLettersSorted) != 7 {
		return errors.New("letters parameter " + uniqueLettersSorted + " must be exactly 7 unique characters. it was " + strconv.Itoa(len(uniqueLettersSorted)) + ".")
	}

	if !strings.Contains(uniqueLettersSorted, string(mandatoryChar)) {
		return errors.New("mandatoryChar parameter must be one of the characters in letters string. " + string(mandatoryChar) + " is not in " + uniqueLettersSorted)
	}

	return nil
}

/*
 s is an arbitrary string
 returns a string containing each unique character of s sorted alfabetically and converted to lowecase.
 e.g. uniqueLettersSorted("Hello") -> "ehlo"
*/
func uniqueLettersSorted(s string) string {

	s = strings.ToLower(s)

	//map for unique letters
	m := make(map[string]bool)

	for i := 0; i < len(s); i++ {
		m[string(s[i])] = true
	}

	//array for letters as strings
	keys := make([]string, len(m))

	i := 0

	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	//	fmt.Println("--->", keys)

	var sb strings.Builder

	//convert strings of letters to 1 string
	for _, c := range keys {
		sb.WriteString(c)
	}
	return sb.String()
}

package main

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	// "encoding/json"
	// "fmt"
	// "log"
	// "net/http"
	// "os"
	// "regexp"
)

const (
	layoutISO_DT = "2006-01-02T15:04:05"
	layoutISO_D  = "2006-01-02"
)

var err error
var httpResp *http.Response
var silentOps bool

func main() {
	filePath, err := findFilePath(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	silentOps = findSilentOps(os.Args)

	list := readListFile(filePath)
	length := len(list)

	for listIndex := range list {
		b := getMetadata(list[listIndex])
		defer httpResp.Body.Close()

		parentFolder := createParentFolder(b.Graphql.ShortcodeMedia.Owner.Username)

		postDate := getPostDate(list[listIndex])
		postID := b.Graphql.ShortcodeMedia.Shortcode
		childFolder := createChildFolder(parentFolder, postDate, postID)

		postType := getPostType(b)
		if postType == "GraphSidecar" {
			saveSidecar(b.Graphql.ShortcodeMedia.Sidecar, childFolder, postID)
			logCurrStatus(listIndex, length)
			continue
		}
		if postType == "GraphImage" {
			displayResources := b.Graphql.ShortcodeMedia.DisplayResources
			dimensions := b.Graphql.ShortcodeMedia.Dimensions
			for index := range displayResources {
				if displayResources[index].ConfigHeight == dimensions.Height && displayResources[index].ConfigWidth == dimensions.Width {
					saveImage(displayResources[index].Src, childFolder, postID)
					logCurrStatus(listIndex, length)
					continue
				}
			}
			continue
		}
		if postType == "GraphVideo" {
			saveVideo(b.Graphql.ShortcodeMedia.VideoURL, childFolder, postID)
			logCurrStatus(listIndex, length)
			continue
		}
	}
	log.Println("Download finished. Closing app now...")
}

func printLog(logMessage string) {
	if !silentOps {
		log.Println(logMessage)
	}
}

func saveSidecar(sidecar Sidecar, savePath string, postID string) {
	edges := sidecar.Edges
	for index := range edges {
		node := edges[index].Node
		typename := node.NodeTypename
		sidecarID := postID + "-" + strconv.Itoa(index+1)
		if typename == "GraphImage" {
			displayResources := node.DisplayResources
			dimensions := node.NodeDimensions
			for index := range displayResources {
				if displayResources[index].ConfigHeight == dimensions.Height && displayResources[index].ConfigWidth == dimensions.Width {
					saveImage(displayResources[index].Src, savePath, sidecarID)
					continue
				}
			}
			continue
		}
		if typename == "GraphVideo" {
			saveVideo(node.VideoURL, savePath, sidecarID)
			continue
		}
	}
}

func saveImage(url string, savePath string, postID string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	saveFile := savePath + "/" + postID + ".jpg"
	imageFile, err := os.Create(saveFile)
	if err != nil {
		log.Fatal(err)
	}
	defer imageFile.Close()

	_, err = io.Copy(imageFile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	printLog(string(postID + ".jpg was saved at folder " + savePath + "."))
}

func saveVideo(url string, savePath string, postID string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	saveFile := savePath + "/" + postID + ".mp4"
	imageFile, err := os.Create(saveFile)
	if err != nil {
		log.Fatal(err)
	}
	defer imageFile.Close()

	_, err = io.Copy(imageFile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	printLog(string(postID + ".mp4 was saved at folder " + savePath + "."))
}

func createParentFolder(username string) string { // (b Body) string {
	parentFolder := string(username) // + "(" + string(b.Graphql.ShortcodeMedia.Owner.DisplayName) + ")"
	if _, err := os.Stat(parentFolder); os.IsNotExist(err) {
		err = os.Mkdir(parentFolder, 0755)
		if err != nil {
			log.Fatal(err)
		}
		printLog(string("Created Folder \"" + parentFolder + "\"."))
	}
	return parentFolder
}

func createChildFolder(parentFolder string, datetime string, postID string) string {
	datetimeObject, err := time.Parse(layoutISO_DT, datetime)
	if err != nil {
		log.Fatal(err)
	}
	date := string(datetimeObject.Format(layoutISO_D))
	childFolder := parentFolder + "/" + date + "__" + postID
	//log.Println(childFolder)
	if _, err := os.Stat(childFolder); os.IsNotExist(err) {
		err = os.Mkdir(childFolder, 0755)
		if err != nil {
			log.Fatal(err)
		}
		printLog(string("Created SubFolder \"" + date + "__" + postID + "\"."))
	}
	//dateObject, err := time.Parse(layoutUS, origDate)
	return childFolder
	//return string(dateObject.Format(layoutISO))
}

func readListFile(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var list []string
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		list = append(list, fileScanner.Text())
	}

	return list
}

func getMetadata(url string) (b Body) {
	resp := getJSONData(url)
	b = newBody(resp)

	return b
}

func getJSONData(url string) *json.Decoder {
	url = url + "?__a=1"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	//request.Header.Set("User-Agent", "PostmanRuntime/7.28.0")

	httpResp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//defer response.Body.Close()

	jsonResp := json.NewDecoder(httpResp.Body)

	return jsonResp
}

func getPostType(b Body) string {
	return b.Graphql.ShortcodeMedia.Typename
}

func getPostDate(url string) string {
	var bodyString string

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		//defer bodyBytes.Close()

		bodyString = string(bodyBytes)
		//log.Info(bodyString)
		//fmt.Println(bodyString)
		regexMatchString := " {\"@context\"(.*)"

		regex, err := regexp.Compile(regexMatchString)
		if err != nil {
			log.Fatal(err)
		}

		bodyString = regex.FindString(bodyString)
		//fmt.Println(regex.FindString(bodyString))
	}

	var pageInfo PageInfo

	json.Unmarshal([]byte(bodyString), &pageInfo)

	//fmt.Println(pageInfo.Datetime)
	//origDate := regex.FindStringSubmatch(b.Graphql.ShortcodeMedia.Description)[1]
	//dateObject, err := time.Parse(layoutUS, origDate)
	return pageInfo.Datetime
	//return string(dateObject.Format(layoutISO))
}

func newBody(response *json.Decoder) (b Body) {
	err = response.Decode(&b)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func logCurrStatus(listIndex int, length int) {
	log.Println(string("Finished " + strconv.Itoa(listIndex+1) + " out of " + strconv.Itoa(length) + "."))
	time.Sleep(5 * time.Second)
}

func findFilePath(args []string) (string, error) {
	for index := range args {
		if strings.Contains(args[index], ".txt") {
			return args[index], nil
		}
	}
	return "", err
}

func findSilentOps(args []string) bool {
	for index := range args {
		if strings.Contains(args[index], "--quiet") || strings.Contains(args[index], "-q") {
			return true
		}
	}
	return false
}

type PageInfo struct {
	Datetime string `json:"uploadDate,omitempty"`
}

// 	link := "https://www.instagram.com/p/CMcLaT8FezH/?__a=1"

// 	req, err := http.NewRequest("GET", link, nil)
// 	if err != nil {
// 		// handle err
// 	}
// 	req.Header.Set("User-Agent", "PostmanRuntime/7.28.0")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		// handle err
// 	}
// 	defer resp.Body.Close()

// 	// x, err := ioutil.ReadAll(req.Body)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// fmt.Println(string(x))

// 	content := json.NewDecoder(resp.Body)

// 	var b Body
// 	err = content.Decode(&b)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	//fmt.Println(b.Graphql.ShortcodeMedia.Description)
// 	//return

// 	regexMatchString := "^Photo by (.*) on (.*)\\. (.*)$"

// 	regex, err := regexp.Compile(regexMatchString)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	fmt.Println(regex.FindStringSubmatch(b.Graphql.ShortcodeMedia.Description)[1])
// 	fmt.Println(regex.FindStringSubmatch(b.Graphql.ShortcodeMedia.Description)[2])

// 	err = os.Mkdir(regex.FindStringSubmatch(b.Graphql.ShortcodeMedia.Description)[1], 0755)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// text := "Photo by IZ*ONE 아이즈원 on March 15, 2021. May be a black-and-white image of 1 person and flower."

// 	// r, _ := regexp.Compile("^Photo by (.*) on (.*)\\. (.*)$")

// 	// fmt.Println(r.FindStringSubmatch(text)[2])

// 	return
// }

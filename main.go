package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-ego/gse"
)

type segRequest struct {
	Text string
}

type segResponse struct {
	Words []string
}

var (
	seger gse.Segmenter
)

func main() {
	fmt.Println("[+] TextSeg start")
	seger.LoadDict("touhou.txt,touhou2.txt,networds.txt,chs.txt,cht.txt,jpn.txt")
	fmt.Println("[+] Done loading dict")
	http.HandleFunc("/s/", segTextSearch)
	http.HandleFunc("/b/", segTextBigram)
	http.HandleFunc("/i/", segTextIndex)
	http.HandleFunc("/d/", segTextDisplay)
	http.HandleFunc("/t/", segTextTouhou)
	http.HandleFunc("/cb/", segCharBigram)
	fmt.Println("[+] Serving ...")
	http.ListenAndServe("0.0.0.0:5005", nil)
}

func unique(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func segTextDisplay(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var requestBodyStr = strings.ToLower(string(requestBody))
	var requestBodyStrNoSpace = strings.Replace(requestBodyStr, " ", "", -1)
	requestBodyStrNoSpace = strings.Replace(requestBodyStrNoSpace, "　", "", -1)
	requestBodyStr = requestBodyStr + " " + requestBodyStrNoSpace

	//splitText := strings.FieldsFunc(requestBodyStr, split)
	//for i := 0; i < len(splitText); i++ {
	//tb := []byte(splitText[i])
	tb := []byte(requestBodyStr)
	segments := seger.Segment(tb)

	// Handle word segmentation results
	// Support for normal mode and search mode two participle,
	// see the comments in the code ToString function.
	// The search mode is mainly used to provide search engines
	// with as many keywords as possible
	//fmt.Println(gse.ToString(segments, true))
	ret := [][]string{}
	for i := 0; i < len(segments); i++ {
		seg := segments[i]
		//wordType := seg.Token().Pos()
		//if wordType == "TH名詞" {
		ele := []string{seg.Token().Text(), seg.Token().Pos()}
		ret = append(ret, ele)
		//fmt.Printf("%s=>%s\n", )
		//}
	}
	//}

	js, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func segTextTouhou(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var requestBodyStr = strings.ToLower(string(requestBody))
	var requestBodyStrNoSpace = strings.Replace(requestBodyStr, " ", "", -1)
	requestBodyStrNoSpace = strings.Replace(requestBodyStrNoSpace, "　", "", -1)
	requestBodyStr = requestBodyStr + " " + requestBodyStrNoSpace

	//splitText := strings.FieldsFunc(requestBodyStr, split)
	//for i := 0; i < len(splitText); i++ {
	//tb := []byte(splitText[i])
	tb := []byte(requestBodyStr)
	segments := seger.Segment(tb)

	// Handle word segmentation results
	// Support for normal mode and search mode two participle,
	// see the comments in the code ToString function.
	// The search mode is mainly used to provide search engines
	// with as many keywords as possible
	//fmt.Println(gse.ToString(segments, true))
	ret := []string{}
	for i := 0; i < len(segments); i++ {
		seg := segments[i]
		wordType := seg.Token().Pos()
		if wordType == "TH名詞" {
			ret = append(ret, seg.Token().Text())
		}
	}

	uniqueRet := unique(ret)
	//}

	js, err := json.Marshal(uniqueRet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func segTextSearch(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var requestBodyStr = strings.ToLower(string(requestBody))
	var requestBodyStrNoSpace = strings.Replace(requestBodyStr, " ", "", -1)
	requestBodyStrNoSpace = strings.Replace(requestBodyStrNoSpace, "　", "", -1)
	requestBodyStr = requestBodyStr + " " + requestBodyStrNoSpace

	segResult := make([]string, 0, len(requestBodyStr))
	splitText := strings.FieldsFunc(requestBodyStr, split)
	for i := 0; i < len(splitText); i++ {
		tmpSegResult := seger.Cut(splitText[i], true)
		segResult = append(segResult, tmpSegResult...)
	}
	resp := segResponse{segResult}
	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func segCharBigram(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var requestBodyStr = strings.ToLower(string(requestBody))
	var requestBodyStrNoSpace = strings.Replace(requestBodyStr, " ", "", -1)
	requestBodyStrNoSpace = strings.Replace(requestBodyStrNoSpace, "　", "", -1)
	requestBodyStr = requestBodyStr + " " + requestBodyStrNoSpace

	chars := make([]rune, 0, len(requestBodyStr))
	splitText := strings.FieldsFunc(requestBodyStr, split_space)
	for i := 0; i < len(splitText); i++ {
		chars = append(chars, []rune(splitText[i])...)
	}
	bigrams := make([]string, 0, len(chars)-1)
	for j := 0; j < len(chars)-1; j++ {
		bigrams = append(bigrams, string([]rune{chars[j], chars[j+1]}))
	}
	resp := segResponse{bigrams}
	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func segTextBigram(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var requestBodyStr = strings.ToLower(string(requestBody))
	var requestBodyStrNoSpace = strings.Replace(requestBodyStr, " ", "", -1)
	requestBodyStrNoSpace = strings.Replace(requestBodyStrNoSpace, "　", "", -1)
	requestBodyStr = requestBodyStr + " " + requestBodyStrNoSpace

	segResult := make([]string, 0, len(requestBodyStr))
	splitText := strings.FieldsFunc(requestBodyStr, split)
	for i := 0; i < len(splitText); i++ {
		tmpSegResult := seger.Cut(splitText[i], true)
		segResult = append(segResult, tmpSegResult...)
	}
	bigrams := make([]string, 0, len(segResult)-1)
	for j := 0; j < len(segResult)-1; j++ {
		bigrams = append(bigrams, segResult[j]+segResult[j+1])
	}
	segResult = append(segResult, bigrams...)
	resp := segResponse{segResult}
	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func segTextIndex(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var requestBodyStr = strings.ToLower(string(requestBody))
	var requestBodyStrNoSpace = strings.Replace(requestBodyStr, " ", "", -1)
	requestBodyStrNoSpace = strings.Replace(requestBodyStrNoSpace, "　", "", -1)
	requestBodyStr = requestBodyStr + " " + requestBodyStrNoSpace

	segResult := make([]string, 0, len(requestBodyStr))
	splitText := strings.FieldsFunc(requestBodyStr, split)
	for i := 0; i < len(splitText); i++ {
		tmpSegResult := seger.CutSearch(splitText[i], true)
		segResult = append(segResult, tmpSegResult...)
	}

	allWords := unique(segResult)
	resp := segResponse{allWords}
	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func split_space(r rune) bool {
	return r == '\t' ||
	r == ' ' ||
	r == '　' ||
	r == '\n' ||
	r == '\r' ||
	r == '\v' ||
	r == '\f' ||
	r == '\x00' ||
	r == '\x84' ||
	r == '\x85' ||
	r == '\x91' ||
	r == '\x80' ||
	r == '\xa0' ||
	r == '\x8d' ||
	r == '\x97' ||
	r == '\x9e' ||
	r == '\x95' ||
	r == '\x8e' ||
	r == '\x9a' ||
	r == '\x96' ||
	r == '\x01' ||
	r == '\x02' ||
	r == '\x03' ||
	r == '\x04' ||
	r == '\x05' ||
	r == '\x06' ||
	r == '\x07' ||
	r == '\x08' ||
	r == '\x09' ||
	r == '\x0a' ||
	r == '\x0b' ||
	r == '\x0c' ||
	r == '\x0d' ||
	r == '\x0e' ||
	r == '\x0f' ||
	r == '\x10' ||
	r == '\x11' ||
	r == '\x12' ||
	r == '\x13' ||
	r == '\x14' ||
	r == '\x15' ||
	r == '\x16' ||
	r == '\x17' ||
	r == '\x18' ||
	r == '\x19' ||
	r == '\x1a' ||
	r == '\x1b' ||
	r == '\x1c' ||
	r == '\x1d' ||
	r == '\x1e' ||
	r == '\x1f' ||
	r == '…' ||
	r == '♡';
}

func split(r rune) bool {
	return r == ':' ||
		r == '.' ||
		r == '\n' ||
		r == '\r' ||
		r == '[' ||
		r == ']' ||
		r == ' ' ||
		r == '\t' ||
		r == '\v' ||
		r == '\f' ||
		r == '{' ||
		r == '}' ||
		r == '-' ||
		r == '_' ||
		r == '=' ||
		r == '+' ||
		r == '`' ||
		r == '~' ||
		r == '!' ||
		r == '@' ||
		r == '#' ||
		r == '$' ||
		r == '%' ||
		r == '^' ||
		r == '&' ||
		r == '*' ||
		r == '(' ||
		r == ')' ||
		r == ';' ||
		r == '\'' ||
		r == '"' ||
		r == ',' ||
		r == '<' ||
		r == '>' ||
		r == '/' ||
		r == '?' ||
		r == '\\' ||
		r == '|' ||
		r == '－' ||
		r == '＞' ||
		r == '＜' ||
		r == '。' ||
		r == '，' ||
		r == '《' ||
		r == '》' ||
		r == '【' ||
		r == '】' ||
		r == '　' ||
		r == '？' ||
		r == '！' ||
		r == '￥' ||
		r == '…' ||
		r == '（' ||
		r == '）' ||
		r == '、' ||
		r == '：' ||
		r == '；' ||
		r == '·' ||
		r == '「' ||
		r == '」' ||
		r == '『' ||
		r == '』' ||
		r == '〔' ||
		r == '〕' ||
		r == '［' ||
		r == '］' ||
		r == '｛' ||
		r == '｝' ||
		r == '｟' ||
		r == '｠' ||
		r == '〉' ||
		r == '〈' ||
		r == '〖' ||
		r == '〗' ||
		r == '〘' ||
		r == '〙' ||
		r == '〚' ||
		r == '〛' ||
		r == '゠' ||
		r == '＝' ||
		r == '‥' ||
		r == '※' ||
		r == '＊' ||
		r == '〽' ||
		r == '〓' ||
		r == '〇' ||
		r == '＂' ||
		r == '“' ||
		r == '”' ||
		r == '‘' ||
		r == '’' ||
		r == '＃' ||
		r == '＄' ||
		r == '％' ||
		r == '＆' ||
		r == '＇' ||
		r == '＋' ||
		r == '．' ||
		r == '／' ||
		r == '＠' ||
		r == '＼' ||
		r == '＾' ||
		r == '＿' ||
		r == '｀' ||
		r == '｜' ||
		r == '～' ||
		r == '｡' ||
		r == '｢' ||
		r == '｣' ||
		r == '､' ||
		r == '･' ||
		r == 'ｰ' ||
		r == 'ﾟ' ||
		r == '￠' ||
		r == '￡' ||
		r == '￢' ||
		r == '￣' ||
		r == '￤' ||
		r == '￨' ||
		r == '￩' ||
		r == '￪' ||
		r == '￫' ||
		r == '￬' ||
		r == '￭' ||
		r == '￮' ||
		r == '・' ||
		r == '◊' ||
		r == '→' ||
		r == '←' ||
		r == '↑' ||
		r == '↓' ||
		r == '↔' ||
		r == '—'
}

package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

type GoTransClient struct {
    TranslateAPIUrl *url.URL
    UrlValues       *url.Values
}

func (g *GoTransClient) BuildUrlParams(sl, tl, text string) {
    val := url.Values{}
    val.Set("client", "p")
    val.Set("ie", "UTF-8")
    val.Set("oe", "UTF-8")
    val.Set("sl", sl)
    val.Set("tl", tl)
    val.Set("text", text)
    g.UrlValues = &val
}

func (g *GoTransClient) ApiUrlToString() string {
    apiUrl := url.URL{}
    apiUrl.Scheme = "http"
    apiUrl.Host = "translate.google.com"
    apiUrl.Path = "translate_a/t"
    g.TranslateAPIUrl = &apiUrl

    if g.UrlValues != nil {
        g.TranslateAPIUrl.RawQuery = g.UrlValues.Encode()
    }
    return g.TranslateAPIUrl.String()
}

func (g *GoTransClient) NewApiUrl() {
    apiUrl := url.URL{}
    apiUrl.Scheme = "http"
    apiUrl.Host = "translate.google.com"
    apiUrl.Path = "translate_a/t"
    // return &apiUrl
}

func (g *GoTransClient) GetTranslateResp(sl, tl, text string) []byte {
    // sl := "en"
    // tl := "zh-CN"
    // text := "test"

    g.BuildUrlParams(sl, tl, text)
    urlStr := g.ApiUrlToString()

    req, _ := http.NewRequest("GET", urlStr, nil)
    req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64)")
    client := &http.Client{}
    resp, _ := client.Do(req)
    body, _ := ioutil.ReadAll(resp.Body)
    return body
}

func NewApiUrl() *url.URL {
    apiUrl := url.URL{}
    apiUrl.Scheme = "http"
    apiUrl.Host = "translate.google.com"
    apiUrl.Path = "translate_a/t"
    return &apiUrl
}

type sentences struct {
    Trans        string `json:"trans"`
    Orig         string `json:"orig"`
    Translit     string `json:"translit"`
    Src_translit string `json:"src_translit"`
}

type dict struct {
    Pos       string   `json:"pos"`
    Terms     []string `json:"terms"`
    Base_form string   `json:"base_form"`
    Pos_enum  int8     `json:"pos_enum"`
    Entry     []entry
}

type entry struct {
    Word                string   `json:"word"`
    Reverse_translation []string `json:"reverse_translation"`
    Score               float32  `json:"score"`
}

type TranslateResult struct {
    Sentences []sentences `json:"sentences"`
    Dict      []dict      `json:"dict"`
}

func PrettyResponse(b []byte) {
    var tr TranslateResult
    err := json.Unmarshal(b, &tr)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }

    if len(tr.Dict) < 1 {
        for i := 0; i < len(tr.Sentences); i++ {
            fmt.Printf("[sentences]:\t%s\n", tr.Sentences[i].Trans)
        }
    }

    for i := 0; i < len(tr.Dict); i++ {
        d := tr.Dict[i]
        // terms := tr.Dict[i].Terms
        // fmt.Printf("[%s]\t%s\n", tr.Dict[i].Pos, strings.Join(terms, ","))
        fmt.Printf("%s [%s] \n", d.Base_form, d.Pos)
        for j := 0; j < len(d.Entry); j++ {
            e := d.Entry[j]
            fmt.Printf("\t%s\t%s\n", e.Word, strings.Join(e.Reverse_translation, ","))
        }
    }

}

func main() {
    var sl, tl string
    flag.StringVar(&sl, "f", "en", "source_language")    //source_language
    flag.StringVar(&tl, "t", "zh-CN", "target_language") //target_language
    flag.Parse()
    text := strings.Join(flag.Args(), " ")
    g := GoTransClient{}
    b := g.GetTranslateResp(sl, tl, text)
    PrettyResponse(b)
}

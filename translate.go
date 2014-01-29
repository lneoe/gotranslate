package main

import (
    // "errors"
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

// build url parameter
// 't' client will receiver non-standard json format
// change client to something other than 't' to get standard json response
// "client=p&ie=utf-8&oe=utf-8&sl=en&tl=zh-CN&text=test"
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

// return the url string
// http://translate.google.com/translate_a/t?client=p&ie=utf-8&oe=utf-8&sl=en&tl=zh-CN&text=test
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

// get the translate result
// return []byte
func (g *GoTransClient) GetTranslateResp(sl, tl, text string) (b []byte, err error) {
    // sl := "en"
    // tl := "zh-CN"
    // text := "test"

    g.BuildUrlParams(sl, tl, text)
    urlStr := g.ApiUrlToString()

    req, _ := http.NewRequest("GET", urlStr, nil)
    req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64)")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        // fmt.Printf("Error: %v\n", err)
        return nil, err
    } else {
        // body, err := ioutil.ReadAll(resp.Body)
        // if err != nil {
        //     fmt.Printf("Error: %v\n", err)
        //     return nil, err
        // }
        body, _ := ioutil.ReadAll(resp.Body)
        return body, nil
    }

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

type spell struct {
    Spell_html_res  string `json:"spell_html_res"`
    Spell_res       string `json:"spell_res"`
    Correction_type []int8 `json:"correction_type"`
    Related         bool   `json:"related"`
}

type TranslateResult struct {
    Sentences []sentences `json:"sentences"`
    Dict      []dict      `json:"dict"`
    Src       string      `json:"src"`
    Spell     spell       `json:"spell"`
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
    if tr.Spell.Spell_res != "" {
        fmt.Println("Did you mean: %s", tr.Spell.Spell_res)
    }

}

func main() {
    var sl, tl string
    var v, help bool
    flag.StringVar(&sl, "f", "en", "source_language,default 'en'")       //source_language
    flag.StringVar(&tl, "t", "zh-CN", "target_language,default 'zh-CN'") //target_language
    flag.BoolVar(&v, "v", false, "reverse 'f' and 'v' options")          //reverse the "f" and "v"
    flag.BoolVar(&help, "help", false, "for help")
    flag.Parse()
    if help || len(flag.Args()) == 0 {
        fmt.Println("Usage:\n  translate [-f <from language>] [-t <to language>] [-v] <data>")
        flag.PrintDefaults()
        fmt.Println("eg.\n  translate -f en -t zh-CN test\n  translate -v test")
    } else {
        text := strings.Join(flag.Args(), " ")
        if v {
            sl, tl = tl, sl
        }

        g := GoTransClient{}
        b, err := g.GetTranslateResp(sl, tl, text)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
        } else {
            PrettyResponse(b)
        }
    }

}

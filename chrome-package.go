package main

import (
    "log"
    "regexp"
    "strings"
    "strconv"
    "net/http"
    "io/ioutil"
    "encoding/json"
)

type Version struct {
    X86   string `json:"x86"`
    X64   string `json:"x64"`
    Mac   string `json:"mac"`
    Appid string `json:"appid"`
}

type Response struct {
    Urls    []string `json:"urls"`
    Version string   `json:"version"`
    Sha256    string   `json:"sha256"`
    Size    int      `json:"size"`
}

type Archs struct {
    X86 Response `json:"x86"`
    X64 Response `json:"x64"`
    Mac Response `json:"mac"`
}

type Chrome struct {
    Stable Archs `json:"stable"`
    Beta   Archs `json:"beta"`
    Dev    Archs `json:"dev"`
    Canary Archs `json:"canary"`
}

func readFile(filename string) []byte {
    data, err := ioutil.ReadFile(filename)

    if err != nil {
        log.Fatal(err)
    }

    return data
}

func writeFile(filename string, data []byte) {
    err := ioutil.WriteFile(filename, data, 0644)

    if err != nil {
        log.Fatal(err)
    }
}

func info() map[string]Version {
    var info map[string]Version

    data := readFile("./info.json")
    err := json.Unmarshal(data, &info)

    if err != nil {
        log.Fatal(err)
    }

    return info
}

func httpDo(url string, data string) []byte {
    client := &http.Client{}

    req, err := http.NewRequest("POST", url, strings.NewReader(data))
    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.Do(req)

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    return body
}

func request(ver, arch string) []byte {
    req := string(readFile("./request.xml"))
    info := info()

    if arch == "mac" {
        req = strings.Replace(req, "{ap}", info[ver].Mac, -1)
    } else if arch == "x86" {
        req = strings.Replace(req, "{ap}", info[ver].X86, -1)
    } else {
        req = strings.Replace(req, "{ap}", info[ver].X64, -1)
    }
    if arch == "mac" {
        if ver == "canary" {
            req = strings.Replace(req, "{appid}", "com.google.Chrome.Canary", -1)
        } else {
            req = strings.Replace(req, "{appid}", "com.google.Chrome", -1)
        }
        req = strings.Replace(req, "{platform}", "mac", -1)
        req = strings.Replace(req, "{version}", "46.0.2490.86", -1)
        req = strings.Replace(req, "{arch}", "x64", -1)
    } else {
        req = strings.Replace(req, "{appid}", info[ver].Appid, -1)
        req = strings.Replace(req, "{platform}", "win", -1)
        req = strings.Replace(req, "{version}", "6.3", -1)
        req = strings.Replace(req, "{arch}", arch, -1)
    }

    log.Print(req)

    return httpDo("https://tools.google.com/service/update2", req)
}

func baseUrl(data []byte) []string {
    var result []string

    re, _ := regexp.Compile("codebase=\"(https://[^\\s]*)\"")
    res := re.FindAllSubmatch(data, -1)

    for _, v := range res {
        result = append(result, string(v[1]))
    }

    return result
}

func installerFilename(data []byte) string {
    var result string

    re, _ := regexp.Compile("run=\"([0-9a-zA-z._-]+)\"")
    res := re.FindSubmatch(data)

    result = string(res[1])

    return result
}

func version(data []byte) string {
    var result string

    re, _ := regexp.Compile("manifest version=\"([0-9.]+)\"")
    res := re.FindSubmatch(data)

    result = string(res[1])

    return result
}

func sha256(data []byte) string {
    var result string

    re, _ := regexp.Compile("hash_sha256=\"([0-9a-z]+)\"")
    res := re.FindSubmatch(data)

    result = string(res[1])

    return result
}

func size(data []byte) int {
    var result int

    re, _ := regexp.Compile("size=\"([0-9]+)\"")
    res := re.FindSubmatch(data)

    result, _ = strconv.Atoi(string(res[1]))

    return result
}

func parse(data []byte) Response {
    var response Response

    urls := baseUrl(data)
    installer := installerFilename(data)

    for _, url := range urls {
        response.Urls = append(response.Urls, url + installer)
    }

    response.Version = version(data)
    response.Sha256 = sha256(data)
    response.Size = size(data)

    return response
}

func chrome() {
    var chromeList Chrome

    chromeList.Stable.X86 = parse(request("stable", "x86"))
    chromeList.Beta.X86 = parse(request("beta", "x86"))
    chromeList.Dev.X86 = parse(request("dev", "x86"))
    chromeList.Canary.X86 = parse(request("canary", "x86"))

    chromeList.Stable.X64 = parse(request("stable", "x64"))
    chromeList.Beta.X64 = parse(request("beta", "x64"))
    chromeList.Dev.X64 = parse(request("dev", "x64"))
    chromeList.Canary.X64 = parse(request("canary", "x64"))

    chromeList.Stable.Mac = parse(request("stable", "mac"))
    chromeList.Beta.Mac = parse(request("beta", "mac"))
    chromeList.Dev.Mac = parse(request("dev", "mac"))
    chromeList.Canary.Mac = parse(request("canary", "mac"))

    chromeJson, err := json.Marshal(chromeList)
    if err != nil {
        log.Fatal(err)
    }

    writeFile("./chrome.json", chromeJson)

    log.Print("OK")
}

func main() {
    chrome()
}

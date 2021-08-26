package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var (
	domain     string
	configFile string
)

type West struct {
	DomainName string `json:"domain"`
	ICPCode    string `json:"icp"`
	ICPStatus  string `json:"icpstatus"`
}

func md5encode(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

func requestURI() (uri string) {
	account := viper.GetString("account")
	key := viper.GetString("key")
	fmt.Println(account)
	fmt.Println(key)
	// MD5 Hash
	var hash_data string = account + key + "domainname"
	sig := md5encode(hash_data)
	rawCmd := fmt.Sprintf("domainname\r\ncheck\r\nentityname:icp\r\ndomains:%s\r\n.\r\n", domain)
	// URL Encoding
	strCmd := url.QueryEscape(rawCmd)
	return fmt.Sprintf(`http://api.west263.com/api/?userid=%s&strCmd=%s&versig=%s`, account, strCmd, sig)
}

func httpPOST() (content []byte, err error) {
	var client = &http.Client{}
	uri := requestURI()
	data := strings.NewReader(``)
	req, err := http.NewRequest("POST", uri, data)
	if err != nil {
		fmt.Println("Resquest error.")
		fmt.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Response error.")
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	// Convert GBK to UTF-8
	reader := simplifiedchinese.GB18030.NewDecoder().Reader(resp.Body)
	content, err = ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("Content error.")
		fmt.Println(err)
		return nil, err
	}
	return
}

func check() string {
	body, err := httpPOST()
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	// Find String
	re, _ := regexp.Compile("{.*}")
	match := fmt.Sprintln(re.FindString(string(body)))
	// Parse Json
	var icp West
	json.Unmarshal([]byte(match), &icp)
	return icp.ICPStatus
}

func readConf() {
	viper.SetConfigName(configFile)
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func cliAction() cli.ActionFunc {
	return func(c *cli.Context) error {
		if c.NumFlags() == 2 || c.NumFlags() == 4 {
			readConf()
			fmt.Println(domain+":", check())
			return nil
		} else {
			cli.ShowAppHelpAndExit(c, 1)
			return nil
		}
	}
}

func main() {
	app := &cli.App{
		Name:            "icp",
		Usage:           "Check ICP status of domain",
		UsageText:       "icp [Global Options] argument",
		HelpName:        "help",
		HideHelp:        false,
		HideHelpCommand: false,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load configuration from `FILE`",
				Value:       "env",
				Destination: &configFile,
			},
			&cli.StringFlag{
				Name:        "domain",
				Aliases:     []string{"d"},
				Usage:       "Domain name",
				Destination: &domain,
			},
		},
		Action: cliAction(),
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

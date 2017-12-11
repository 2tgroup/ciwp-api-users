package findcountry

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"

	"bitbucket.org/2tgroup/ciwp-api-users/dbconnects"
)

const (
	defaultLink = "https://raw.githubusercontent.com/pirsquare/country-mapper/master/files/country_info.csv"
	defaultFile = "./common/findcountry/country_info.csv"
	collection  = "country"
)

//Country to get start get list
var Country *CountryInfoClient

//TypeLoad  is can load by file,link,database
var TypeLoad string = "database"

func init() {
	client, err := LoadCountryBy(TypeLoad)
	if err != nil {
		panic(err)
	}
	client.ListData()
	client.ListMapByAlpha2()
	Country = client
	log.Info("Loaded common/findcountry")
}

type CountryInfoClient struct {
	Data         []*CountryInfo
	List         []CountryInfo
	ListByAlpha2 []string `json:"list_by_alpha2,omitempty"`
}

func (c *CountryInfoClient) MapByName(name string) *CountryInfo {
	for _, row := range c.Data {
		// check Name field
		if strings.ToLower(row.Name) == strings.ToLower(name) {
			return row
		}

		// check AlternateNames field
		if stringInSlice(strings.ToLower(name), row.AlternateNamesLower()) {
			return row
		}
	}
	return nil
}

func (c *CountryInfoClient) MapByAlpha2(alpha2 string) *CountryInfo {
	for _, row := range c.Data {
		if strings.ToLower(row.Alpha2) == strings.ToLower(alpha2) {
			return row
		}
	}
	return nil
}

func (c *CountryInfoClient) MapByAlpha3(alpha3 string) *CountryInfo {
	for _, row := range c.Data {
		if strings.ToLower(row.Alpha3) == strings.ToLower(alpha3) {
			return row
		}
	}
	return nil
}

func (c *CountryInfoClient) MapByCurrency(currency string) []*CountryInfo {
	rowList := []*CountryInfo{}
	for _, row := range c.Data {
		if stringInSlice(strings.ToLower(currency), row.CurrencyLower()) {
			rowList = append(rowList, row)
		}
	}
	return rowList
}

func (c *CountryInfoClient) MapByCallingCode(callingCode string) []*CountryInfo {
	rowList := []*CountryInfo{}
	for _, row := range c.Data {
		if stringInSlice(strings.ToLower(callingCode), row.CallingCodeLower()) {
			rowList = append(rowList, row)
		}
	}
	return rowList
}

func (c *CountryInfoClient) MapByRegion(region string) []*CountryInfo {
	rowList := []*CountryInfo{}
	for _, row := range c.Data {
		if strings.ToLower(row.Region) == strings.ToLower(region) {
			rowList = append(rowList, row)
		}
	}
	return rowList
}

func (c *CountryInfoClient) MapBySubregion(subregion string) []*CountryInfo {
	rowList := []*CountryInfo{}
	for _, row := range c.Data {
		if strings.ToLower(row.Subregion) == strings.ToLower(subregion) {
			rowList = append(rowList, row)
		}
	}
	return rowList
}

func (c *CountryInfoClient) ListData() {
	if len(c.List) < 1 {
		for _, row := range c.Data {
			c.List = append(c.List, *row)
		}
	}
}

func (c *CountryInfoClient) ListMapByAlpha2() {
	for _, row := range c.Data {
		c.ListByAlpha2 = append(c.ListByAlpha2, (row.Alpha2))
	}
}

type CountryInfo struct {
	Name           string   `json:"Name,omitempty" bson:"Name"`
	AlternateNames []string `json:"AlternateNames,omitempty" bson:"AlternateNames"`
	Alpha2         string   `json:"Alpha2,omitempty" bson:"Alpha2"`
	Alpha3         string   `json:"Alpha3,omitempty" bson:"Alpha3"`
	Currency       []string `json:"Currency,omitempty" bson:"Currency"`
	CallingCode    []string `json:"CallingCode,omitempty" bson:"CallingCode"`
	Region         string   `json:"Region,omitempty" bson:"Region"`
	Subregion      string   `json:"Subregion,omitempty" bson:"Subregion"`
}

func (c *CountryInfo) AlternateNamesLower() []string {
	updated := []string{}
	for _, alternateName := range c.AlternateNames {
		updated = append(updated, strings.ToLower(alternateName))
	}
	return updated
}

func (c *CountryInfo) CurrencyLower() []string {
	updated := []string{}
	for _, currency := range c.Currency {
		updated = append(updated, strings.ToLower(currency))
	}
	return updated
}

func (c *CountryInfo) CallingCodeLower() []string {
	updated := []string{}
	for _, callingCode := range c.CallingCode {
		updated = append(updated, strings.ToLower(callingCode))
	}
	return updated
}

func readCSVFromURL(fileURL string) ([][]string, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ';'
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func readCSVFromLocal(filePathName string) ([][]string, error) {

	filepath, err := filepath.Abs(filePathName)
	if err != nil {
		fmt.Println("Read file local:", err)
	}
	csvFile, err := os.Open(filepath)
	// automatically call Close() at the end of current method
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	reader.Comma = ';'
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Pass in an optional url if you would like to use your own downloadable csv file for country's data.
// This is useful if you prefer to host the data file yourself or if you have modified some of the fields
// for your specific use case.
func Load(specifiedURL ...string) (*CountryInfoClient, error) {
	var fileURL string
	var data [][]string

	if len(specifiedURL) > 0 {
		fileURL = specifiedURL[0]
	} else {
		fileURL = defaultFile
	}
	//checking it url or file path
	u, err := url.Parse(fileURL)
	if err != nil {

		log.Fatal(err)
	}

	if u.Host == "" {
		data, err = readCSVFromLocal(fileURL)
	} else {
		data, err = readCSVFromURL(fileURL)
	}

	if err != nil {
		return nil, err
	}

	recordList := []*CountryInfo{}
	for idx, row := range data {
		// skip header
		if idx == 0 {
			continue
		}

		// get name
		name := strings.Split(row[0], ",")[:1][0]

		// use commonly used & altSpellings names as AlternateNames
		alternateNames := strings.Split(row[0], ",")[1:]
		alternateNames = append(alternateNames, strings.Split(row[8], ",")...)

		record := &CountryInfo{
			Name:           name,
			AlternateNames: alternateNames,
			Alpha2:         row[2],
			Alpha3:         row[4],
			Currency:       strings.Split(row[5], ","),
			CallingCode:    strings.Split(row[6], ","),
			Region:         row[10],
			Subregion:      row[11],
		}
		//condition := dbconnect.MongodbToBson(record)
		countRow := dbconnect.CountRowsInCollection(collection, record)
		if countRow < 1 {
			dbconnect.InserToCollection(collection, record)
		}
		recordList = append(recordList, record)
	}
	return &CountryInfoClient{Data: recordList}, nil
}

func runMigration() {
	var data [][]string
	data, _ = readCSVFromLocal(defaultFile)
	for idx, row := range data {
		// skip header
		if idx == 0 {
			continue
		}
		// get name
		name := strings.Split(row[0], ",")[:1][0]
		// use commonly used & altSpellings names as AlternateNames
		alternateNames := strings.Split(row[0], ",")[1:]
		alternateNames = append(alternateNames, strings.Split(row[8], ",")...)

		record := &CountryInfo{
			Name:           name,
			AlternateNames: alternateNames,
			Alpha2:         row[2],
			Alpha3:         row[4],
			Currency:       strings.Split(row[5], ","),
			CallingCode:    strings.Split(row[6], ","),
			Region:         row[10],
			Subregion:      row[11],
		}
		//condition := dbconnect.MongodbToBson(record)
		countRow := dbconnect.CountRowsInCollection(collection, record)
		if countRow < 1 {
			dbconnect.InserToCollection(collection, record)
		}
	}
}

func LoadCountryBy(from string) (*CountryInfoClient, error) {
	switch from {
	case "link":
		return Load([]string{defaultLink}...)
	case "file":
		return Load()
	default:
		runMigration()
		return LoadFromDatabase()
	}
}

func LoadFromDatabase() (*CountryInfoClient, error) {

	recordList := []*CountryInfo{}

	dataFound, err := dbconnect.SearchDataInCollection(collection, echo.Map{}, 0, 0)

	byteData, errMar := json.Marshal(dataFound)

	if errMar != nil {
		return nil, err
	}
	err = json.Unmarshal(byteData, &recordList)

	return &CountryInfoClient{Data: recordList}, err
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

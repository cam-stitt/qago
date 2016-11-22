package qago_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/cam-stitt/qago"
	"github.com/fatih/color"
	"github.com/sclevine/agouti"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

const tick = "\u2705"
const cross = "\u2717"

var directory string
var noColor bool

var bold = color.New(color.Bold).SprintFunc()
var boldBlue = color.New(color.FgBlue, color.Bold).SprintfFunc()
var boldYellow = color.New(color.FgYellow, color.Bold).SprintfFunc()
var boldGreen = color.New(color.FgGreen, color.Bold).SprintfFunc()
var boldRed = color.New(color.FgRed, color.Bold).SprintfFunc()
var boldCyan = color.New(color.FgCyan, color.Bold).SprintfFunc()

func init() {
	flag.StringVar(&directory, "case-dir", "./fixtures", "the directory for the suites")
	flag.BoolVar(&noColor, "no-color", false, "Disable color output")

	flag.Parse()

	if noColor {
		color.NoColor = true // disables colorized output
	}
}

type SeleniumTestSuite struct {
	suite.Suite
	FileList []os.FileInfo
	Driver   *agouti.WebDriver
	Page     *agouti.Page
}

func (sts *SeleniumTestSuite) SetupSuite() {
	sts.Driver = agouti.PhantomJS()
	//	sts.Driver = agouti.ChromeDriver()
	sts.Driver.Start()
}

func (sts *SeleniumTestSuite) getSelection(predicate *qago.Predicate) *agouti.Selection {
	t := sts.T()

	var output string
	if predicate.Name != "" {
		output = predicate.Name
	} else {
		output = fmt.Sprintf("Type: %s, Selector: %s", predicate.Type, predicate.Selector)
	}
	fmt.Printf("%s %s - ", boldYellow("Predicate:"), output)

	page := sts.Page
	selectorType := predicate.Type
	var selectorPrefix string
	if predicate.Multi {
		selectorPrefix = qago.Multi
	} else if predicate.First {
		selectorPrefix = qago.First
	} else {
		selectorPrefix = qago.Find
	}

	args := qago.ExtractMethodArgs(selectorType, predicate)
	method := qago.FindSelectMethod(selectorType, selectorPrefix)

	pageValue := reflect.ValueOf(page)
	methodValue := pageValue.MethodByName(method)

	returnValues := methodValue.Call(args)

	if len(returnValues) != 1 {
		fmt.Println(boldRed("Not Found"))
		t.Fatalf("Expected a single return value from method: %s", method)
	}

	selection := returnValues[0].Interface().(*agouti.Selection)

	fmt.Println(boldGreen("Found"))

	return selection
}

func (sts *SeleniumTestSuite) runAction(selection *agouti.Selection, actions []qago.Action) {
	var err error
	var output string
	for idx, action := range actions {
		output = boldYellow("Action %d:", idx)
		if action.Name != "" {
			output = fmt.Sprintf("%s %s", output, action.Name)
		} else {
			output = fmt.Sprintf("%s %s", action.Type)
		}
		fmt.Printf("%s - ", output)

		switch action.Type {
		case qago.Click:
			err = selection.Click()
		case qago.Fill:
			err = selection.Fill(action.Text)
		case qago.Check:
			err = selection.Check()
		case qago.Clear:
			err = selection.Clear()
		case qago.DoubleClick:
			err = selection.DoubleClick()
		case qago.SendKeys:
			err = selection.SendKeys(action.Text)
		}
		if !sts.NoError(err) {
			fmt.Println(boldRed("Failure"))
		}

		fmt.Println(boldGreen("Success"))

		if action.Wait != "" {
			wait, err := time.ParseDuration(action.Wait)
			sts.NoError(err)
			time.Sleep(wait)
		}
		sts.runAssertions(selection, action.Assertions)
	}
}

func (sts *SeleniumTestSuite) runAssertions(selectable interface{}, assertions []qago.Assertion) {
	page := sts.Page

	var output string
	for idx, assertion := range assertions {
		output = boldYellow("Assertion %d:", idx)
		if assertion.Name != "" {
			output = fmt.Sprintf("%s %s", output, assertion.Name)
		}
		fmt.Printf("%s - ", output)

		if len(assertion.Query) > 0 {
			currentURL, err := page.URL()
			sts.NoError(err)
			actualURL, err := url.Parse(currentURL)
			sts.NoError(err)
			query := actualURL.Query()
			for _, arg := range assertion.Query {
				actual := query.Get(arg.Key)
				sts.Equal(arg.Value, actual)
			}
		}
		if assertion.URL != "" {
			currentURL, err := page.URL()
			sts.NoError(err)
			sts.Equal(assertion.URL, currentURL)
		}
		if selectable != nil {
			selection, ok := selectable.(*agouti.Selection)
			sts.True(ok)
			if assertion.Text != "" {
				text, err := selection.Text()
				if !sts.NoError(err) && !sts.Equal(assertion.Text, text) {
					fmt.Println(boldRed("Failure"))
				}
			}
			for _, attribute := range assertion.Attributes {
				actual, err := selection.Attribute(attribute.Key)
				if !sts.NoError(err) && !sts.Equal(attribute.Value, actual) {
					fmt.Println(boldRed("Failure"))
				}
			}
		}

		fmt.Println(boldGreen("Success"))
	}
	fmt.Print("\n")
}

func (sts *SeleniumTestSuite) runTestCase(testCase qago.Case) {
	capabilities := agouti.NewCapabilities().Browser(testCase.Browser)
	page, err := sts.Driver.NewPage(agouti.Desired(capabilities))
	sts.NoError(err)
	sts.Page = page

	err = page.Navigate(testCase.Location)
	sts.NoError(err)

	prefix := "Case:"
	suffix := testCase.Name
	fmt.Printf("%s %s\n", boldCyan(prefix), bold(suffix))
	for i := 0; i < len(fmt.Sprintf("%s %s", prefix, suffix)); i++ {
		fmt.Print("=")
	}
	fmt.Print("\n")

	for idx, step := range testCase.Steps {
		prefix = fmt.Sprintf("Step %d:", idx)
		suffix = step.Name
		fmt.Printf("%s %s\n", boldBlue(prefix), bold(suffix))
		for i := 0; i < len(fmt.Sprintf("%s %s", prefix, suffix)); i++ {
			fmt.Print("-")
		}
		fmt.Print("\n")

		selection := sts.getSelection(step.Predicate)

		sts.runAction(selection, step.Actions)

		if step.Wait != "" {
			wait, err := time.ParseDuration(step.Wait)
			sts.NoError(err)
			time.Sleep(wait)
		}
	}

	fmt.Println(bold("Running global assertions"))
	fmt.Println("-------------------------")
	sts.runAssertions(nil, testCase.Assertions)
}

func (sts *SeleniumTestSuite) TestSeleniumSuite() {
	for _, file := range sts.FileList {
		data, err := ioutil.ReadFile(filepath.Join(directory, file.Name()))
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		testCase := qago.Case{}
		err = yaml.Unmarshal(data, &testCase)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		sts.runTestCase(testCase)
	}
}

func (sts *SeleniumTestSuite) TearDownSuite() {
	sts.Driver.Stop()
}

func TestAutomatedSuite(t *testing.T) {
	var fileList []os.FileInfo
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		matched, err := regexp.MatchString(".+.yml", fileName)
		if err != nil {
			panic(err)
		}
		if matched {
			fileList = append(fileList, file)
		}
	}

	suite.Run(t, &SeleniumTestSuite{
		FileList: fileList,
	})
}

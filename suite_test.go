package qago_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"reflect"
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

var bold = color.New(color.Bold).SprintFunc()
var boldBlue = color.New(color.FgBlue, color.Bold).SprintFunc()
var boldYellow = color.New(color.FgYellow, color.Bold).SprintFunc()
var boldGreen = color.New(color.FgGreen, color.Bold).SprintFunc()
var boldRed = color.New(color.FgRed, color.Bold).SprintFunc()

func init() {
	flag.StringVar(&directory, "testdir", "", "the directory for the suites")
}

type SeleniumTestSuite struct {
	suite.Suite
	TestCase *qago.Case
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
		if action.Name != "" {
			output = action.Name
		} else {
			output = fmt.Sprintf("%d", idx)
		}
		fmt.Printf("%s %s - ", boldYellow("Action:"), output)

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
		sts.runAssertions(selection, action.Assertions)
	}
}

func (sts *SeleniumTestSuite) runAssertions(selectable interface{}, assertions []qago.Assertion) {
	page := sts.Page

	var output string
	for idx, assertion := range assertions {
		if assertion.Name != "" {
			output = assertion.Name
		} else {
			output = fmt.Sprintf("%d", idx)
		}
		fmt.Printf("%s %s - ", boldYellow("Assertion:"), output)

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
		if selectable != nil {
			selection, ok := selectable.(*agouti.Selection)
			sts.True(ok)
			if assertion.Text != "" {
				text, err := selection.Text()
				if !sts.NoError(err) || !sts.Equal(assertion.Text, text) {
					fmt.Println(boldRed("Failure"))
				}
			}
			for _, attribute := range assertion.Attributes {
				actual, err := selection.Attribute(attribute.Key)
				if !sts.NoError(err) || !sts.Equal(attribute.Value, actual) {
					fmt.Println(boldRed("Failure"))
				}
			}
		}

		fmt.Println(boldGreen("Success"))
	}
	fmt.Print("\n")
}

func (sts *SeleniumTestSuite) TestSeleniumSuite() {
	testCase := sts.TestCase

	capabilities := agouti.NewCapabilities().Browser(testCase.Browser)
	page, err := sts.Driver.NewPage(agouti.Desired(capabilities))
	sts.NoError(err)
	sts.Page = page

	err = page.Navigate(testCase.Location)
	sts.NoError(err)

	fmt.Printf("%s\n\n", bold(testCase.Name))

	var output string
	for idx, step := range testCase.Steps {
		output = boldBlue("Step ", idx)
		if step.Name != "" {
			output = fmt.Sprintf("%s: %s", output, bold(step.Name))
		}
		characterLen := len(output)
		fmt.Printf("%s\n", output)
		for i := 0; i < characterLen; i++ {
			fmt.Print("=")
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

	fmt.Println("Running global assertions")
	fmt.Println("=========================")
	sts.runAssertions(nil, testCase.Assertions)
}

func (sts *SeleniumTestSuite) TearDownSuite() {
	sts.Driver.Stop()
}

func TestAutomatedSuite(t *testing.T) {
	data, err := ioutil.ReadFile("fixtures/helloworld.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	testCase := qago.Case{}
	err = yaml.Unmarshal(data, &testCase)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	suite.Run(t, &SeleniumTestSuite{
		TestCase: &testCase,
	})
}

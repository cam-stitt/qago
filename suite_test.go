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
	"github.com/sclevine/agouti"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

const tick = "\u2705"
const cross = "\u2717"

var directory string

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
	suffix := cross

	defer func() {
		fmt.Printf("%s: %s\n", output, suffix)
	}()

	if predicate.Name != "" {
		output = predicate.Name
	} else {
		output = fmt.Sprintf("Type: %s, Selector: %s", predicate.Type, predicate.Selector)
	}

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
		t.Fatalf("Expected a single return value from method: %s", method)
	}

	selection := returnValues[0].Interface().(*agouti.Selection)

	suffix = tick
	return selection
}

func (sts *SeleniumTestSuite) runAction(selection *agouti.Selection, actions []qago.Action) {
	var err error
	for idx, action := range actions {
		if action.Name != "" {
			fmt.Println(action.Name)
		} else {
			fmt.Printf("Action %d\n", idx)
		}

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
		sts.NoError(err)

		sts.runAssertions(selection, action.Assertions)
	}
}

func (sts *SeleniumTestSuite) runAssertions(selectable interface{}, assertions []qago.Assertion) {
	page := sts.Page

	var output string
	suffix := cross

	for idx, assertion := range assertions {
		defer func() {
			fmt.Printf("%s: %s\n", output, suffix)
		}()

		if assertion.Name != "" {
			output = assertion.Name
		} else {
			output = fmt.Sprintf("Assertion %d", idx)
		}

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
				sts.NoError(err)
				sts.Equal(assertion.Text, text)
			}
			for _, attribute := range assertion.Attributes {
				actual, err := selection.Attribute(attribute.Key)
				sts.NoError(err)
				sts.Equal(attribute.Value, actual)
			}
		}

		suffix = tick
	}
}

func (sts *SeleniumTestSuite) TestSeleniumSuite() {
	testCase := sts.TestCase

	capabilities := agouti.NewCapabilities().Browser(testCase.Browser)
	page, err := sts.Driver.NewPage(agouti.Desired(capabilities))
	sts.NoError(err)
	sts.Page = page

	err = page.Navigate(testCase.Location)
	sts.NoError(err)

	var output string
	for idx, step := range testCase.Steps {
		fmt.Println("=============================")
		output = fmt.Sprintf("Step %d", idx)
		if step.Name != "" {
			output = fmt.Sprintf("%s: %s", output, step.Name)
		}
		fmt.Printf("%s\n", output)
		fmt.Println("=============================")

		selection := sts.getSelection(step.Predicate)

		sts.runAction(selection, step.Actions)

		if step.Wait != "" {
			wait, err := time.ParseDuration(step.Wait)
			sts.NoError(err)
			time.Sleep(wait)
		}
	}

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

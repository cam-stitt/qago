package qago_test

import (
	"io/ioutil"
	"log"
	"testing"
	"flag"
	"net/url"
	"time"

	"github.com/cam-stitt/qago"
	"gopkg.in/yaml.v2"
	"github.com/stretchr/testify/suite"
	"github.com/sclevine/agouti"
)

var directory string
func init() {
	flag.StringVar(&directory, "testdir", "", "the directory for the suites")
}

type SeleniumTestSuite struct {
	suite.Suite
	SeleniumSuite *qago.Suite
	Driver        *agouti.WebDriver
	Page          *agouti.Page
}

func (sts *SeleniumTestSuite) SetupSuite() {
	sts.Driver = agouti.PhantomJS()
//	sts.Driver = agouti.ChromeDriver()
	sts.Driver.Start()
}

func (sts *SeleniumTestSuite) getSelection(predicate *qago.Predicate) *agouti.Selection {
	page := sts.Page
	selector := predicate.Selector
	
	var selection *agouti.Selection
	switch predicate.Type {
	case qago.Find:
		selection = page.Find(selector)
	case qago.FindByButton:
		selection = page.FindByButton(selector)
	case qago.FindByClass:
		selection = page.FindByClass(selector)
	case qago.FindByID:
		selection = page.FindByID(selector)
	case qago.FindByLabel:
		selection = page.FindByLabel(selector)
	case qago.FindByLink:
		selection = page.FindByLink(selector)
	case qago.FindByName:
		selection = page.FindByName(selector)
	case qago.FindByXPath:
		selection = page.FindByXPath(selector)
	case qago.FindForAppium:
		selection = page.FindForAppium(selector, predicate.Text)
	case qago.First:
		selection = page.First(selector)
	case qago.FirstByButton:
		selection = page.FirstByButton(selector)
	case qago.FirstByClass:
		selection = page.FirstByClass(selector)
	case qago.FirstByLabel:
		selection = page.FirstByLabel(selector)
	case qago.FirstByLink:
		selection = page.FirstByLink(selector)
	case qago.FirstByName:
		selection = page.FirstByName(selector)
	case qago.FirstByXPath:
		selection = page.FirstByXPath(selector)
	default:
		sts.T().Fatal("Predicate not found")
	}

	return selection
}

func (sts *SeleniumTestSuite) runAction(selection *agouti.Selection, action *qago.Action) {
	var err error
	switch action.Type {
	case qago.Click:
		err = selection.Click()
	case qago.Fill:
		err = selection.Fill(action.Text)
	case qago.SendKeys:
		err = selection.SendKeys(action.Text)
	}
	sts.NoError(err)
}

func (sts *SeleniumTestSuite) runAssertions(selectable interface{}, assertions []qago.Assertion) {
	for _, assertion := range action.Assertions {
		if len(assertion.Query) > 0 {
			currentURL, err := sts.Page.URL()
			sts.NoError(err)
			actualURL, err := url.Parse(currentURL)
			sts.NoError(err)
			query := actualURL.Query()
			for _, arg := range assertion.Query {
				actual := query.Get(arg.Arg)
				sts.Equal(arg.Value, actual)
			}
		}
		if assertion.Text != "" {
			text, err := selection.Text()
			sts.NoError(err)
			sts.Equal(assertion.Text, text)
		}
		for _, attribute := range assertion.Attributes {
			actual, err := selection.Attribute(attribute.Name)
			sts.NoError(err)
			sts.Equal(attribute.Value, actual)
		}
	}
}

func (sts *SeleniumTestSuite) TestSeleniumSuite() {
	t := sts.T()
	suite := sts.SeleniumSuite

	capabilities := agouti.NewCapabilities().Browser(suite.Browser)
	page, err := sts.Driver.NewPage(agouti.Desired(capabilities))
	sts.NoError(err)

	err = page.Navigate(suite.Location)
	sts.NoError(err)

	for _, step := range *suite.Steps {
		t.Log(step.Name)
		selection := sts.getSelection(step.Predicate)

		for _, action := range *step.Actions {
			t.Log(action.Name)
			sts.runAction(selection, &action)
			sts.runAssertions(selection, action.Assertions)
		}

		if step.Wait != "" {
			wait, err := time.ParseDuration(step.Wait)
			sts.NoError(err)
			time.Sleep(wait)
		}
	}

	sts.runAssertions(nil, suite.Assertions)
	for _, assertion := range *suite.Assertions {
		if len(assertion.Query) > 0 {
			currentURL, err := page.URL()
			sts.NoError(err)
			actualURL, err := url.Parse(currentURL)
			sts.NoError(err)
			query := actualURL.Query()
			for _, arg := range assertion.Query {
				actual := query.Get(arg.Arg)
				sts.Equal(arg.Value, actual)
			}
		}
	}
}

func (sts *SeleniumTestSuite) TearDownSuite() {
	sts.Driver.Stop()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	data, err := ioutil.ReadFile("fixtures/helloworld.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	seleniumSuite := qago.Suite{}
	err = yaml.Unmarshal(data, &seleniumSuite)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	suite.Run(t, &SeleniumTestSuite{
		SeleniumSuite: &seleniumSuite,
	})
}

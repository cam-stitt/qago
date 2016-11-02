package qago_test

import (
	"io/ioutil"
	"log"
	"testing"
	"flag"

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
	Driver *agouti.WebDriver
}

func (sts *SeleniumTestSuite) SetupSuite() {
	sts.Driver = agouti.ChromeDriver()
	sts.Driver.Start()
}

func (sts *SeleniumTestSuite) getSelection(page *agouti.Page, predicate *qago.Predicate) *agouti.Selection {
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
	t := sts.T()
	var err error
	switch action.Type {
	case qago.Click:
		err = selection.Click()
		if err != nil {
			t.Fatal(err)
		}
	case qago.Fill:
		err = selection.Fill(action.Text)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func (sts *SeleniumTestSuite) TestSeleniumSuite() {
	suite := sts.SeleniumSuite
	t := sts.T()
	capabilities := agouti.NewCapabilities().Browser(suite.Browser)
	page, err := sts.Driver.NewPage(agouti.Desired(capabilities))
	if err != nil {
		t.Fatal(err)
	}
	err = page.Navigate(suite.Location)
	if err != nil {
		t.Fatal(err)
	}
	for _, step := range *suite.Steps {
		t.Log(step.Name)
		selection := sts.getSelection(page, step.Predicate)

		for _, action := range *step.Actions {
			sts.runAction(selection, &action)
		}
		t.Log(selection)
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

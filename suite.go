package qago

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"time"

	yaml "gopkg.in/yaml.v1"

	"github.com/fatih/color"
	"github.com/sclevine/agouti"
)

var bold = color.New(color.Bold).SprintFunc()
var boldBlue = color.New(color.FgBlue, color.Bold).SprintfFunc()
var boldYellow = color.New(color.FgYellow, color.Bold).SprintfFunc()
var boldGreen = color.New(color.FgGreen, color.Bold).SprintfFunc()
var boldRed = color.New(color.FgRed, color.Bold).SprintfFunc()
var boldCyan = color.New(color.FgCyan, color.Bold).SprintfFunc()

type Suite struct {
	NoColor   bool
	Directory string
	FileList  []os.FileInfo
	Driver    *agouti.WebDriver
	Page      *agouti.Page
}

func (s *Suite) readFiles() error {
	fileList := []os.FileInfo{}

	files, err := ioutil.ReadDir(s.Directory)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		matched, err := regexp.MatchString(".+.yml", fileName)
		if err != nil {
			return err
		}
		if matched {
			fileList = append(fileList, file)
		}
	}

	s.FileList = fileList

	return nil
}

func (s *Suite) Run() int {
	defer func(begin time.Time) {
		fmt.Printf("Took: %s\n", time.Since(begin))
	}(time.Now())

	if s.NoColor {
		color.NoColor = true
	}

	s.Driver = agouti.PhantomJS()
	//	sts.Driver = agouti.ChromeDriver()
	s.Driver.Start()

	err := s.readFiles()
	if err != nil {
		return 1
	}

	for _, file := range s.FileList {
		data, err := ioutil.ReadFile(filepath.Join(s.Directory, file.Name()))
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		testCase := Case{}
		err = yaml.Unmarshal(data, &testCase)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		err = s.runTestCase(testCase)
		if err != nil {
			return 1
		}
	}

	return 0
}

func (s *Suite) runTestCase(testCase Case) error {
	capabilities := agouti.NewCapabilities().Browser(testCase.Browser)
	page, err := s.Driver.NewPage(agouti.Desired(capabilities))
	if err != nil {
		return err
	}

	s.Page = page

	err = page.Navigate(testCase.Location)
	if err != nil {
		return err
	}

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

		selection, err := s.getSelection(step.Predicate)
		if err != nil {
			return err
		}

		err = s.runAction(selection, step.Actions)
		if err != nil {
			return err
		}

		if step.Wait != "" {
			wait, err := time.ParseDuration(step.Wait)
			if err != nil {
				return err
			}
			time.Sleep(wait)
		}
	}

	fmt.Println(bold("Running global assertions"))
	fmt.Println("-------------------------")
	return s.runAssertions(nil, testCase.Assertions)
}

func (s *Suite) runAssertions(selectable interface{}, assertions []Assertion) error {
	page := s.Page

	var output string
	for idx, assertion := range assertions {
		output = boldYellow("Assertion %d:", idx)
		if assertion.Name != "" {
			output = fmt.Sprintf("%s %s", output, assertion.Name)
		}
		fmt.Printf("%s - ", output)

		if len(assertion.Query) > 0 {
			currentURL, err := page.URL()
			if err != nil {
				return err
			}
			actualURL, err := url.Parse(currentURL)
			if err != nil {
				return err
			}
			query := actualURL.Query()
			for _, arg := range assertion.Query {
				actual := query.Get(arg.Key)
				if arg.Value != actual {
					return fmt.Errorf("Assertion failed: Expected %s, got %s", arg.Value, actual)
				}
			}
		}
		if assertion.URL != "" {
			currentURL, err := page.URL()
			if err != nil {
				return err
			}
			if assertion.URL != currentURL {
				return fmt.Errorf("Assertion failed: Expected %s, got %s", assertion.URL, currentURL)
			}
		}
		if selectable != nil {
			selection, ok := selectable.(*agouti.Selection)
			if !ok {
				return fmt.Errorf("Failed to get selection")
			}
			if assertion.Text != "" {
				text, err := selection.Text()
				if err != nil || assertion.Text != text {
					fmt.Println(boldRed("Failure\n"))
					return fmt.Errorf("Failure")
				}
			}
			for _, attribute := range assertion.Attributes {
				actual, err := selection.Attribute(attribute.Key)
				if err != nil || attribute.Value != actual {
					fmt.Println(boldRed("Failure\n"))
					return fmt.Errorf("Failure")
				}
			}
		}

		fmt.Println(boldGreen("Success"))
	}
	fmt.Print("\n")

	return nil
}

func (s *Suite) getSelection(predicate *Predicate) (*agouti.Selection, error) {
	var output string
	if predicate.Name != "" {
		output = predicate.Name
	} else {
		output = fmt.Sprintf("Type: %s, Selector: %s", predicate.Type, predicate.Selector)
	}
	fmt.Printf("%s %s - ", boldYellow("Predicate:"), output)

	page := s.Page
	selectorType := predicate.Type
	var selectorPrefix string
	if predicate.Multi {
		selectorPrefix = Multi
	} else if predicate.First {
		selectorPrefix = First
	} else {
		selectorPrefix = Find
	}

	args := ExtractMethodArgs(selectorType, predicate)
	method := FindSelectMethod(selectorType, selectorPrefix)

	pageValue := reflect.ValueOf(page)
	methodValue := pageValue.MethodByName(method)

	returnValues := methodValue.Call(args)

	if len(returnValues) != 1 {
		fmt.Println(boldRed("Not Found"))
		return nil, fmt.Errorf("Expected a single return value from method: %s", method)
	}

	selection := returnValues[0].Interface().(*agouti.Selection)

	fmt.Println(boldGreen("Found"))

	return selection, nil
}

func (s *Suite) runAction(selection *agouti.Selection, actions []Action) error {
	var err error
	var output string
	for idx, action := range actions {
		output = boldYellow("Action %d:", idx)
		if action.Name != "" {
			output = fmt.Sprintf("%s %s", output, action.Name)
		} else {
			output = fmt.Sprintf("%s %s", output, action.Type)
		}
		fmt.Printf("%s - ", output)

		switch action.Type {
		case Click:
			err = selection.Click()
		case Fill:
			err = selection.Fill(action.Text)
		case Check:
			err = selection.Check()
		case Clear:
			err = selection.Clear()
		case DoubleClick:
			err = selection.DoubleClick()
		case SendKeys:
			err = selection.SendKeys(action.Text)
		}
		if err != nil {
			fmt.Println(boldRed("Failure"))
			return fmt.Errorf("Failure")
		}

		fmt.Println(boldGreen("Success"))

		if action.Wait != "" {
			wait, err := time.ParseDuration(action.Wait)
			if err != nil {
				return err
			}
			time.Sleep(wait)
		}
		return s.runAssertions(selection, action.Assertions)
	}

	return nil
}

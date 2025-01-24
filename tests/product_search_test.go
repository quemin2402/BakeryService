package tests

import (
	"github.com/tebeka/selenium"
	"strings"
	"testing"
	"time"
)

func TestProductSearch(t *testing.T) {
	const serverURL = "http://127.0.0.1:4444/wd/hub"

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, serverURL)
	if err != nil {
		t.Fatalf("Error connecting to Selenium Server: %v", err)
	}
	defer wd.Quit()

	t.Log("Session created successfully")

	if err := wd.Get("https://bakeryservice.onrender.com/static/menu.html"); err != nil {
		t.Fatalf("Failed to load page: %v", err)
	}
	t.Log("Page loaded successfully")

	searchInput, err := wd.FindElement(selenium.ByID, "searchInput")
	if err != nil {
		t.Fatalf("Search input not found: %v", err)
	}
	t.Log("Search input found")
	if err := searchInput.SendKeys("Croissant"); err != nil {
		t.Fatalf("Failed to type in search input: %v", err)
	}
	t.Log("Text entered into search input")

	searchButton, err := wd.FindElement(selenium.ByID, "filterButton")
	if err != nil {
		t.Fatalf("Search button not found: %v", err)
	}
	t.Log("Search button found")
	if err := searchButton.Click(); err != nil {
		t.Fatalf("Failed to click search button: %v", err)
	}
	t.Log("Search button clicked")

	err = wd.WaitWithTimeoutAndInterval(
		selenium.Condition(func(wd selenium.WebDriver) (bool, error) {
			elements, err := wd.FindElements(selenium.ByCSSSelector, "img.card-img-top.menu-card-img")
			if err != nil {
				t.Log("Elements not yet available, retrying...")
				return false, nil
			}

			found := 0
			for _, element := range elements {
				altText, err := element.GetAttribute("alt")
				if err == nil && strings.Contains(altText, "Croissant") {
					t.Logf("Found product: %s", altText)
					found++
				}
			}

			if found > 0 {
				t.Logf("Found %d matching products", found)
				return true, nil
			}

			t.Log("No matching products yet, waiting...")
			return false, nil
		}),
		60*time.Second, 500*time.Millisecond,
	)
	if err != nil {
		t.Fatalf("Timeout waiting for product elements: %v", err)
	}
	t.Log("Matching products found successfully.")

}

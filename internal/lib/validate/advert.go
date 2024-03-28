package validate

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"net/http"
	"path/filepath"

	"github.com/rigbyel/ad-market/internal/lib/request"
	"github.com/rigbyel/ad-market/internal/models/constraints"
)

// validates advert according to constraints from models/constraints
func ValidateAdvert(ad request.AdvertRequest) []string {
	errs := []string{}

	// check if advert header exists
	if len(ad.Header) == 0 {
		errs = append(errs, "header is required")
	}

	// check advert header length
	if len(ad.Header) > constraints.AdvertHeaderMaxLen {
		errs = append(errs, "advert header is too long")
	}

	// check advert body length
	if len(ad.Body) > constraints.AdvertBodyMaxLen {
		errs = append(errs, "advert body is too long")
	}

	// check if advert has price
	if ad.Price == 0 {
		errs = append(errs, "price is required")
	}

	// check if price in permmitted range
	if ad.Price > constraints.MaxPrice {
		errs = append(errs, "price is to big")
	}

	if ad.Price < constraints.MinPrice {
		errs = append(errs, "price is too small")
	}

	// validate image
	err := validateImage(ad.ImageURL)
	if err != nil {
		errs = append(errs, err.Error())
	}

	return errs
}

// validates image according to size and extention constraints
func validateImage(imgURL string) error {
	if imgURL == "" {
		return nil
	}

	// get response from image url
	resp, err := http.Get(imgURL)
	if err != nil {
		return fmt.Errorf("invalid image url")
	}
	defer resp.Body.Close()

	// check if image extention is valid
	ext := filepath.Ext(imgURL)
	if _, ok := constraints.ImageExtentions[ext]; !ok {
		return fmt.Errorf("wrong image extention: %s", ext)
	}

	// decode image
	m, _, err := image.Decode(resp.Body)
	if err != nil {
		return fmt.Errorf("invalid image url %w", err)
	}

	// get image bounds
	g := m.Bounds()

	// get height and width
	height := g.Dy()
	width := g.Dx()

	// check if image size is valid
	if height > constraints.ImageMaxHeight || width > constraints.ImageMaxWidth {
		return fmt.Errorf("image is too big")
	}

	if height < constraints.ImageMinHeight || width < constraints.ImageMinWidth {
		return fmt.Errorf("image is too small")
	}

	return nil
}

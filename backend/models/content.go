package models

type SectionVisibility struct {
	Title   string `json:"title" dynamodbav:"title"`
	Visible bool   `json:"visible" dynamodbav:"visible"`
}

type HeroSection struct {
	SectionVisibility
	Subtitle        string `json:"subtitle" dynamodbav:"subtitle"`
	BackgroundImage string `json:"backgroundImage" dynamodbav:"backgroundImage"`
}

type AboutSection struct {
	SectionVisibility
	Description string `json:"description" dynamodbav:"description"`
	Image       string `json:"image" dynamodbav:"image"`
}

type ContactSection struct {
	SectionVisibility
	Address     string `json:"address" dynamodbav:"address"`
	Phone       string `json:"phone" dynamodbav:"phone"`
	Email       string `json:"email" dynamodbav:"email"`
	MapEmbedURL string `json:"mapEmbedUrl" dynamodbav:"mapEmbedUrl"`
}

type SiteContent struct {
	ID      string         `json:"id" dynamodbav:"id"`
	Hero    HeroSection    `json:"hero" dynamodbav:"hero"`
	About   AboutSection   `json:"about" dynamodbav:"about"`
	Contact ContactSection `json:"contact" dynamodbav:"contact"`
}

func DefaultSiteContent() SiteContent {
	return SiteContent{
		ID: "site-content",
		Hero: HeroSection{
			SectionVisibility: SectionVisibility{Title: "Welcome to Our Company", Visible: true},
			Subtitle:          "Delivering excellence since day one",
			BackgroundImage:   "",
		},
		About: AboutSection{
			SectionVisibility: SectionVisibility{Title: "About Us", Visible: true},
			Description:       "We are a passionate team dedicated to providing the best services.",
			Image:             "",
		},
		Contact: ContactSection{
			SectionVisibility: SectionVisibility{Title: "Contact Us", Visible: true},
			Address:           "123 Business Street, City, Country",
			Phone:             "+62 812 3456 7890",
			Email:             "info@company.com",
			MapEmbedURL:       "",
		},
	}
}

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SectionVisibility struct {
	Title   string `json:"title" bson:"title"`
	Visible bool   `json:"visible" bson:"visible"`
}

type HeroSection struct {
	SectionVisibility `bson:",inline"`
	Subtitle          string `json:"subtitle" bson:"subtitle"`
	BackgroundImage   string `json:"backgroundImage" bson:"backgroundImage"`
}

type AboutSection struct {
	SectionVisibility `bson:",inline"`
	Description       string `json:"description" bson:"description"`
	Image             string `json:"image" bson:"image"`
}

type ContactSection struct {
	SectionVisibility `bson:",inline"`
	Address           string `json:"address" bson:"address"`
	Phone             string `json:"phone" bson:"phone"`
	Email             string `json:"email" bson:"email"`
	MapEmbedURL       string `json:"mapEmbedUrl" bson:"mapEmbedUrl"`
}

type SiteContent struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Hero    HeroSection        `json:"hero" bson:"hero"`
	About   AboutSection       `json:"about" bson:"about"`
	Contact ContactSection     `json:"contact" bson:"contact"`
}

func DefaultSiteContent() SiteContent {
	return SiteContent{
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

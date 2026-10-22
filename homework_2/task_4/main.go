package main

import "fmt"

type Profile struct {
	Name       string
	Age        int
	Email      string
	Newsletter bool
}

type ProfilePatch struct {
	Name       *string
	Age        *int
	Email      *string
	Newsletter *bool
}

func ApplyPatch(p *Profile, patch ProfilePatch) {
	if p == nil {
		return
	}
	if patch.Name != nil {
		p.Name = *patch.Name
	}
	if patch.Age != nil {
		p.Age = *patch.Age
	}
	if patch.Email != nil {
		p.Email = *patch.Email
	}
	if patch.Newsletter != nil {
		p.Newsletter = *patch.Newsletter
	}
}

func main() {
	profile := Profile{
		Name:       "Анна",
		Age:        30,
		Email:      "anna@example.com",
		Newsletter: true,
	}

	patch := ProfilePatch{
		Age:        new(31),
		Newsletter: new(bool),
	}

	fmt.Println(profile)
	ApplyPatch(&profile, patch)
	fmt.Println(profile)
}

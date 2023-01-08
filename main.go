package main

import (
	"fmt"
	"github.com/ant0ine/go-webfinger"
	"github.com/go-ap/activitypub"
	"github.com/go-ap/client"
	"github.com/microcosm-cc/bluemonday"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: fediclient @user@example.com")
	}
	email := os.Args[1]
	cleansedEmail := strings.Trim(email, "@")

	fmt.Println(cleansedEmail)

	p := bluemonday.StripTagsPolicy()
	r, err := webfinger.Lookup(cleansedEmail, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(r.Subject)
	for i := range r.Links {
		fmt.Println(r.Links[i].Href)
		if r.Links[i].Rel == "self" {
			c := client.New()
			i, err := c.LoadIRI(activitypub.IRI(r.Links[i].Href))
			if err != nil {
				log.Println(err)
			}
			fmt.Println(i.GetType())

			activitypub.OnActor(i, func(actor *activitypub.Actor) error {
				act := *actor
				fmt.Println(p.Sanitize(act.Summary.String()))

				outbox := act.Outbox.GetLink()
				o, err := c.LoadIRI(outbox)
				if err != nil {
					log.Println(err)
				}
				activitypub.OnOrderedCollection(o, func(out *activitypub.OrderedCollection) error {
					colData, err := c.LoadIRI(out.Last.GetLink())
					if err != nil {
						log.Println(err)
					}
					activitypub.OnOrderedCollectionPage(colData, func(col *activitypub.OrderedCollectionPage) error {
						for item := range col.OrderedItems {
							activitypub.OnActivity(col.OrderedItems[item], func(itemData *activitypub.Activity) error {
								activitypub.OnObject(itemData.Object, func(objectData *activitypub.Object) error {
									fmt.Println("➡️", objectData.Published, "➡️", p.Sanitize(objectData.Content.String()))
									return nil
								})
								return nil
							})
						}
						return nil
					})
					return nil
				})

				return nil
			})

		}
	}

}

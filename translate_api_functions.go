package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

func AuthTranslate(jsonPath, projectID string) (*translate.Client, context.Context, error) {
	ctx := context.Background()
	client, err := translate.NewClient(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		log.Fatal(err)
		return client, ctx, err
	}
	return client, ctx, nil

}

// this is directly copy/pasted from Google example
func translateTextWithModel(targetLanguage, text, model string) (string, error) {

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", fmt.Errorf("language.Parse: %v", err)
	}
	client, ctx, err := AuthTranslate("google-secret.json", "103373479946395174633")
	if err != nil {
		return "", fmt.Errorf("translate.NewClient: %v", err)
	}
	defer client.Close()
	resp, err := client.Translate(ctx, []string{text}, lang, &translate.Options{
		Model: model, // Either "nmt" or "base".
	})
	if err != nil {
		return "", fmt.Errorf("translate: %v", err)
	}
	if len(resp) == 0 {
		return "", nil
	}
	return resp[0].Text, nil
}

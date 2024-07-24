package provider

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var (
	_ function.Function = SesSmtpPasswordV4{}
)

func NewSesSmtpPasswordV4Function() function.Function {
	return SesSmtpPasswordV4{}
}

type SesSmtpPasswordV4 struct{}

func (r SesSmtpPasswordV4) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "ses_smtp_password_v4"
}

func (r SesSmtpPasswordV4) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Get SES SMTP pasword with secret and region",
		MarkdownDescription: "Convert secret access key into an SES SMTP password by applying [AWS's documented Sigv4 conversion algorithm](https://docs.aws.amazon.com/ses/latest/DeveloperGuide/smtp-credentials.html#smtp-credentials-convert).",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "secret",
				MarkdownDescription: "Secert access key",
			},
			function.StringParameter{
				Name:                "region",
				MarkdownDescription: "SES region",
			},
		},
		Return: function.StringReturn{},
	}
}

func hmacSignature(key []byte, value []byte) ([]byte, error) {
	h := hmac.New(sha256.New, key)
	if _, err := h.Write(value); err != nil {
		return []byte(""), err
	}
	return h.Sum(nil), nil
}

func sesSMTPPasswordFromSecretKeySigV4(secret string, region string) (string, error) {
	if secret == "" || region == "" {
		return "", nil
	}
	const version = byte(0x04)
	date := []byte("11111111")
	service := []byte("ses")
	terminal := []byte("aws4_request")
	message := []byte("SendRawEmail")

	rawSig, err := hmacSignature([]byte("AWS4"+secret), date)
	if err != nil {
		return "", err
	}

	if rawSig, err = hmacSignature(rawSig, []byte(region)); err != nil {
		return "", err
	}
	if rawSig, err = hmacSignature(rawSig, service); err != nil {
		return "", err
	}
	if rawSig, err = hmacSignature(rawSig, terminal); err != nil {
		return "", err
	}
	if rawSig, err = hmacSignature(rawSig, message); err != nil {
		return "", err
	}

	versionedSig := make([]byte, 0, len(rawSig)+1)
	versionedSig = append(versionedSig, version)
	versionedSig = append(versionedSig, rawSig...)
	return base64.StdEncoding.EncodeToString(versionedSig), nil
}

func (r SesSmtpPasswordV4) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var secret string
	var region string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &secret, &region))

	if resp.Error != nil {
		return
	}

	ret, err := sesSMTPPasswordFromSecretKeySigV4(secret, region)
	if err != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, ret))
}
